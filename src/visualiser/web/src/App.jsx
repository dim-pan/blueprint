import { useCallback, useEffect, useMemo, useState } from 'react';
import {
  ReactFlow,
  Background,
  Controls,
  MiniMap,
  Panel,
  useNodesState,
  useEdgesState,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import dagre from '@dagrejs/dagre';

const NODE_W = 220;
const NODE_H = 56;
const SUB_NODE_W = 180;
const SUB_NODE_H = 48;

const TYPE_STYLE = {
  component:  { background: '#e6f5df', border: '#3a7a3a' },
  sysreq:     { background: '#e0e6ff', border: '#3344aa' },
  derivedReq: { background: '#fff4d0', border: '#aa7a33' },
  testspec:   { background: '#f5e6ff', border: '#663399' },
  interface:  { background: '#f0f0f0', border: '#666666' },
  gap:        { background: '#ffd7d7', border: '#aa2222' },
};

const EDGE_STYLE = {
  satisfies:    { stroke: '#3a7a3a', strokeWidth: 1.5 },
  depends_on:   { stroke: '#888',    strokeWidth: 1.2, strokeDasharray: '4 4' },
  derived_from: { stroke: '#aa7a33', strokeWidth: 1.2, strokeDasharray: '4 4' },
  allocated_to: { stroke: '#aa7a33', strokeWidth: 1.5 },
  verifies:     { stroke: '#663399', strokeWidth: 1.2 },
};

function truncate(s, n) {
  if (!s) return '';
  return s.length > n ? s.slice(0, n - 1) + '…' : s;
}

function nodeStyle(kind, width = NODE_W) {
  const p = TYPE_STYLE[kind];
  return {
    background: p.background,
    border: `1.5px solid ${p.border}`,
    borderRadius: 6,
    fontSize: 12,
    padding: 8,
    width,
    textAlign: 'center',
    whiteSpace: 'pre-line',
  };
}

function laidOut(nodes, edges, opts = {}) {
  const { w = NODE_W, h = NODE_H, rankdir = 'TB', nodesep = 60, ranksep = 110 } = opts;
  const g = new dagre.graphlib.Graph();
  g.setGraph({ rankdir, nodesep, ranksep, marginx: 20, marginy: 20 });
  g.setDefaultEdgeLabel(() => ({}));
  nodes.forEach((n) => g.setNode(n.id, { width: w, height: h }));
  edges.forEach((e) => g.setEdge(e.source, e.target));
  dagre.layout(g);
  return nodes.map((n) => {
    const { x, y } = g.node(n.id);
    return { ...n, position: { x: x - w / 2, y: y - h / 2 } };
  });
}

// BP-009-VIS-02: Top-level overview — components + system reqs only.
function buildOverview(model) {
  if (!model) return { nodes: [], edges: [] };

  const satisfyingComponents = new Map();
  for (const c of model.components) {
    for (const s of c.Satisfies || []) {
      if (!satisfyingComponents.has(s)) satisfyingComponents.set(s, []);
      satisfyingComponents.get(s).push(c.ID);
    }
  }

  const nodes = [];
  const edges = [];

  for (const c of model.components) {
    nodes.push({
      id: c.ID,
      type: 'default',
      data: { label: `${c.ID}\n${truncate(c.Name, 28)}`, kind: 'component', payload: c },
      style: nodeStyle('component'),
      position: { x: 0, y: 0 },
    });
  }

  for (const r of model.requirements) {
    if (r.DerivedFrom) continue; // system reqs only
    const isUnallocated = !(satisfyingComponents.get(r.ID) || []).length;
    nodes.push({
      id: r.ID,
      type: 'default',
      data: { label: `${r.ID}\n${truncate(r.Title, 30)}`, kind: 'sysreq', payload: r },
      style: nodeStyle(isUnallocated ? 'gap' : 'sysreq'),
      position: { x: 0, y: 0 },
    });
  }

  const ids = new Set(nodes.map((n) => n.id));
  for (const c of model.components) {
    for (const s of c.Satisfies || []) {
      if (ids.has(s)) edges.push({ id: `${c.ID}->sat->${s}`, source: c.ID, target: s, style: EDGE_STYLE.satisfies });
    }
    for (const d of c.DependsOn || []) {
      if (ids.has(d)) edges.push({ id: `${c.ID}->dep->${d}`, source: c.ID, target: d, style: EDGE_STYLE.depends_on });
    }
  }

  return { nodes: laidOut(nodes, edges), edges };
}

// BP-009-VIS-06: Sub-canvas — component + derived reqs + verifying tests + satisfied sysreqs.
function buildComponentSubgraph(model, componentId) {
  const c = model.components.find((x) => x.ID === componentId);
  if (!c) return { nodes: [], edges: [] };

  const derived = model.requirements.filter((r) => r.AllocatedTo === componentId);
  const derivedIds = new Set(derived.map((r) => r.ID));
  const sysReqs = model.requirements.filter((r) => (c.Satisfies || []).includes(r.ID));
  const tests = model.test_specs.filter((t) => (t.Verifies || []).some((v) => derivedIds.has(v)));

  const nodes = [
    {
      id: c.ID,
      type: 'default',
      data: { label: `${c.ID}\n${truncate(c.Name, 24)}`, kind: 'component', payload: c },
      style: nodeStyle('component', SUB_NODE_W),
      position: { x: 0, y: 0 },
    },
    ...sysReqs.map((r) => ({
      id: r.ID,
      type: 'default',
      data: { label: `${r.ID}\n${truncate(r.Title, 26)}`, kind: 'sysreq', payload: r },
      style: nodeStyle('sysreq', SUB_NODE_W),
      position: { x: 0, y: 0 },
    })),
    ...derived.map((r) => ({
      id: r.ID,
      type: 'default',
      data: { label: `${r.ID}\n${truncate(r.Title, 26)}`, kind: 'derivedReq', payload: r },
      style: nodeStyle('derivedReq', SUB_NODE_W),
      position: { x: 0, y: 0 },
    })),
    ...tests.map((t) => ({
      id: t.ID,
      type: 'default',
      data: { label: `${t.ID}\n${truncate(t.Title, 26)}`, kind: 'testspec', payload: t },
      style: nodeStyle('testspec', SUB_NODE_W),
      position: { x: 0, y: 0 },
    })),
  ];

  const edges = [];
  for (const s of c.Satisfies || []) {
    edges.push({ id: `${c.ID}->sat->${s}`, source: c.ID, target: s, style: EDGE_STYLE.satisfies });
  }
  for (const r of derived) {
    edges.push({ id: `${r.ID}->alloc->${c.ID}`, source: r.ID, target: c.ID, style: EDGE_STYLE.allocated_to });
  }
  for (const t of tests) {
    for (const v of t.Verifies || []) {
      if (derivedIds.has(v)) {
        edges.push({ id: `${t.ID}->ver->${v}`, source: t.ID, target: v, style: EDGE_STYLE.verifies });
      }
    }
  }

  return {
    nodes: laidOut(nodes, edges, { w: SUB_NODE_W, h: SUB_NODE_H, rankdir: 'LR', nodesep: 30, ranksep: 70 }),
    edges,
  };
}

function SourceLine({ payload }) {
  if (!payload?.SourceFile) return null;
  return (
    <div style={{ marginTop: 10, fontSize: 11, color: '#777' }}>
      <b>Source:</b> {payload.SourceFile}:{payload.LineNumber}
    </div>
  );
}

function ElementDetails({ node }) {
  const { data } = node;
  const p = data.payload;
  return (
    <div>
      <h3 style={{ margin: 0, fontSize: 14 }}>
        {p.ID || p.Name} <span style={{ color: '#888', fontWeight: 400 }}>({data.kind})</span>
      </h3>
      <p style={{ fontSize: 12, color: '#222', marginTop: 6 }}>{p.Title || p.Name || ''}</p>
      {p.Statement && <Block label="Statement">{p.Statement}</Block>}
      {p.Responsibility && <Block label="Responsibility">{p.Responsibility}</Block>}
      {p.Priority && <Row label="Priority">{p.Priority}</Row>}
      {p.Pattern && <Row label="Pattern">{p.Pattern}</Row>}
      {p.DerivedFrom && <Row label="Derived from">{p.DerivedFrom}</Row>}
      {p.AllocatedTo && <Row label="Allocated to">{p.AllocatedTo}</Row>}
      {p.Satisfies?.length > 0 && <Row label="Satisfies">{p.Satisfies.join(', ')}</Row>}
      {p.DependsOn?.length > 0 && <Row label="Depends on">{p.DependsOn.join(', ')}</Row>}
      {p.Verifies?.length > 0 && <Row label="Verifies">{p.Verifies.join(', ')}</Row>}
      {p.Given && <Block label="Given">{p.Given}</Block>}
      {p.Expect && <Block label="Expect">{p.Expect}</Block>}
      <SourceLine payload={p} />
    </div>
  );
}

function Row({ label, children }) {
  return <div style={{ marginTop: 6, fontSize: 12 }}><b>{label}:</b> {children}</div>;
}
function Block({ label, children }) {
  return (
    <div style={{ marginTop: 10, fontSize: 12, background: '#fafafa', padding: 8, borderRadius: 4 }}>
      <b>{label}:</b><br/>{children}
    </div>
  );
}

// BP-009-VIS-05: List view — system reqs satisfied, then derived reqs with tests nested.
function ComponentList({ model, componentId, onSelectElement }) {
  const c = model.components.find((x) => x.ID === componentId);
  const sysReqs = model.requirements.filter((r) => (c.Satisfies || []).includes(r.ID));
  const derived = model.requirements.filter((r) => r.AllocatedTo === componentId);
  const testsByReq = new Map();
  for (const t of model.test_specs) {
    for (const v of t.Verifies || []) {
      if (!testsByReq.has(v)) testsByReq.set(v, []);
      testsByReq.get(v).push(t);
    }
  }

  const idLink = (id, kind, payload) => (
    <span
      onClick={() => onSelectElement({ id, data: { kind, payload } })}
      style={{ color: '#3344aa', cursor: 'pointer', textDecoration: 'underline' }}
    >{id}</span>
  );

  return (
    <div>
      <h4 style={sectionH}>Responsibility</h4>
      <div style={{ fontSize: 12, color: '#333' }}>{c.Responsibility || <i style={{ color: '#999' }}>none</i>}</div>

      <h4 style={sectionH}>System requirements satisfied</h4>
      {sysReqs.length === 0 ? (
        <div style={emptyStyle}>None.</div>
      ) : sysReqs.map((r) => (
        <div key={r.ID} style={{ marginBottom: 6 }}>
          <div>{idLink(r.ID, 'sysreq', r)} <span style={{ color: '#333' }}>— {r.Title}</span></div>
          {r.SourceFile && <div style={{ fontSize: 10, color: '#888', marginLeft: 10 }}>{r.SourceFile}:{r.LineNumber}</div>}
        </div>
      ))}

      <h4 style={sectionH}>Derived requirements</h4>
      {derived.length === 0 ? (
        <div style={emptyStyle}>This component has no derived requirements yet.</div>
      ) : derived.map((r) => {
        const ts = testsByReq.get(r.ID) || [];
        return (
          <div key={r.ID} style={reqCard}>
            <div>{idLink(r.ID, 'derivedReq', r)} <span style={{ color: '#333' }}>— {r.Title}</span></div>
            {r.Priority && <div style={{ fontSize: 11, color: '#888' }}>{r.Priority}</div>}
            {r.Statement && <div style={{ fontSize: 12, color: '#444', marginTop: 4 }}>{r.Statement}</div>}
            {r.SourceFile && <div style={{ fontSize: 10, color: '#888', marginTop: 4 }}>{r.SourceFile}:{r.LineNumber}</div>}
            {ts.length > 0 && (
              <div style={{ marginTop: 6, paddingLeft: 10, borderLeft: '2px solid #e0d4b8' }}>
                {ts.map((t) => (
                  <div key={t.ID} style={{ fontSize: 11, marginTop: 3 }}>
                    {idLink(t.ID, 'testspec', t)} <span style={{ color: '#555' }}>— {t.Title}</span>
                  </div>
                ))}
              </div>
            )}
          </div>
        );
      })}

      {c.DependsOn?.length > 0 && (
        <>
          <h4 style={sectionH}>Depends on</h4>
          <div style={{ fontSize: 12 }}>{c.DependsOn.map((d, i) => (
            <span key={d}>{i > 0 && ', '}{idLink(d, 'component', model.components.find((x) => x.ID === d) || { ID: d })}</span>
          ))}</div>
        </>
      )}

      <SourceLine payload={c} />
    </div>
  );
}

const sectionH = { margin: '16px 0 6px', fontSize: 12, color: '#555', textTransform: 'uppercase', letterSpacing: 0.5 };
const emptyStyle = { fontSize: 12, color: '#999', fontStyle: 'italic' };
const reqCard = { marginTop: 8, padding: 8, background: '#fffaf0', border: '1px solid #eed9b0', borderRadius: 4 };

function SubCanvas({ model, componentId, onNodeClick }) {
  const { nodes, edges } = useMemo(() => buildComponentSubgraph(model, componentId), [model, componentId]);
  const [n, , onNodesChange] = useNodesState(nodes);
  const [e, , onEdgesChange] = useEdgesState(edges);
  // remount-on-component-change via key in parent handles node/edge refresh.
  return (
    <div style={{ height: '100%', width: '100%' }}>
      <ReactFlow
        nodes={n}
        edges={e}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onNodeClick={(_, node) => onNodeClick(node)}
        fitView
        minZoom={0.2}
        proOptions={{ hideAttribution: true }}
      >
        <Background />
        <Controls showInteractive={false} />
      </ReactFlow>
    </div>
  );
}

function DetailPane({ model, selected, onClose, onSelectElement }) {
  const [view, setView] = useState('list'); // BP-009-VIS-04: list is default
  useEffect(() => { setView('list'); }, [selected?.id]);

  if (!selected) {
    return (
      <div style={{ padding: 16, fontSize: 12, color: '#666' }}>
        <h3 style={{ margin: 0, fontSize: 14, color: '#222' }}>Click a component</h3>
        <p>Click any component node to open its detail pane. Click a system req for quick details.</p>
      </div>
    );
  }

  if (selected.data.kind !== 'component') {
    return (
      <div style={{ padding: 16 }}>
        <CloseButton onClose={onClose} />
        <ElementDetails node={selected} />
      </div>
    );
  }

  const c = selected.data.payload;
  const btn = (active) => ({
    fontSize: 12,
    padding: '4px 10px',
    border: '1px solid #bbb',
    background: active ? '#3344aa' : '#fafafa',
    color: active ? '#fff' : '#222',
    cursor: 'pointer',
    borderRadius: 4,
  });

  return (
    <div style={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <div style={{ padding: '12px 16px', borderBottom: '1px solid #eee', display: 'flex', alignItems: 'center', gap: 10 }}>
        <div style={{ flex: 1 }}>
          <div style={{ fontSize: 14, fontWeight: 600 }}>{c.ID} <span style={{ color: '#888', fontWeight: 400 }}>— {c.Name}</span></div>
          {c.SourceFile && <div style={{ fontSize: 10, color: '#888' }}>{c.SourceFile}:{c.LineNumber}</div>}
        </div>
        <button style={btn(view === 'list')} onClick={() => setView('list')}>List</button>
        <button style={btn(view === 'canvas')} onClick={() => setView('canvas')}>Canvas</button>
        <CloseButton onClose={onClose} />
      </div>
      {view === 'list' ? (
        <div style={{ flex: 1, overflowY: 'auto', padding: 16 }}>
          <ComponentList model={model} componentId={c.ID} onSelectElement={onSelectElement} />
        </div>
      ) : (
        <div style={{ flex: 1 }}>
          <SubCanvas
            key={c.ID}
            model={model}
            componentId={c.ID}
            onNodeClick={(node) => onSelectElement(node)}
          />
        </div>
      )}
    </div>
  );
}

function CloseButton({ onClose }) {
  return (
    <button
      onClick={onClose}
      style={{ fontSize: 14, padding: '2px 8px', border: 'none', background: 'transparent', cursor: 'pointer', color: '#666' }}
      aria-label="Close"
    >✕</button>
  );
}

function Legend() {
  const Row = ({ color, border, label }) => (
    <div style={{ display: 'flex', alignItems: 'center', gap: 6, fontSize: 11 }}>
      <span style={{ display: 'inline-block', width: 14, height: 10, background: color, border: `1px solid ${border}`, borderRadius: 2 }} />
      {label}
    </div>
  );
  return (
    <div style={{ background: '#fff', padding: '8px 10px', border: '1px solid #ddd', borderRadius: 6, display: 'flex', flexDirection: 'column', gap: 4 }}>
      <Row color={TYPE_STYLE.component.background}  border={TYPE_STYLE.component.border}  label="Component" />
      <Row color={TYPE_STYLE.sysreq.background}     border={TYPE_STYLE.sysreq.border}     label="System req" />
      <Row color={TYPE_STYLE.derivedReq.background} border={TYPE_STYLE.derivedReq.border} label="Derived req" />
      <Row color={TYPE_STYLE.testspec.background}   border={TYPE_STYLE.testspec.border}   label="Test spec" />
      <Row color={TYPE_STYLE.gap.background}        border={TYPE_STYLE.gap.border}        label="Unallocated" />
    </div>
  );
}

export default function App() {
  const [model, setModel] = useState(null);
  const [error, setError] = useState(null);
  const [selected, setSelected] = useState(null);

  useEffect(() => {
    fetch('/api/model')
      .then((r) => { if (!r.ok) throw new Error(`HTTP ${r.status}`); return r.json(); })
      .then(setModel)
      .catch((e) => setError(String(e)));
  }, []);

  const overview = useMemo(() => buildOverview(model), [model]);
  const [nodes, setNodes, onNodesChange] = useNodesState(overview.nodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(overview.edges);

  useEffect(() => {
    setNodes(overview.nodes);
    setEdges(overview.edges);
  }, [overview, setNodes, setEdges]);

  const onNodeClick = useCallback((_, node) => setSelected(node), []);
  const clearSelection = useCallback(() => setSelected(null), []);

  if (error) return <div style={{ padding: 20, color: '#a22' }}>Failed to load model: {error}</div>;
  if (!model) return <div style={{ padding: 20 }}>Loading model…</div>;

  const paneOpen = selected !== null;

  return (
    <div style={{ width: '100vw', height: '100vh', display: 'flex' }}>
      <div style={{ flex: 1, position: 'relative', minWidth: 0 }}>
        <ReactFlow
          nodes={nodes}
          edges={edges}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          onNodeClick={onNodeClick}
          onPaneClick={clearSelection}
          fitView
          minZoom={0.1}
        >
          <Background />
          <Controls />
          <MiniMap pannable zoomable />
          <Panel position="bottom-left">
            <Legend />
          </Panel>
        </ReactFlow>
      </div>
      {paneOpen && (
        <div style={{
          width: 460,
          borderLeft: '1px solid #ddd',
          background: '#fff',
          boxShadow: '-2px 0 8px rgba(0,0,0,0.04)',
          display: 'flex',
          flexDirection: 'column',
          overflow: 'hidden',
        }}>
          <DetailPane
            model={model}
            selected={selected}
            onClose={clearSelection}
            onSelectElement={setSelected}
          />
        </div>
      )}
    </div>
  );
}
