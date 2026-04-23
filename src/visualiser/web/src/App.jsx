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

const NODE_W = 200;
const NODE_H = 60;

const TYPE_STYLE = {
  component: { background: '#e6f5df', border: '#3a7a3a' },
  sysreq: { background: '#e0e6ff', border: '#3344aa' },
  derivedReq: { background: '#fff4d0', border: '#aa7a33' },
  testspec: { background: '#f5e6ff', border: '#663399' },
  interface: { background: '#f0f0f0', border: '#666666' },
  gap: { background: '#ffd7d7', border: '#aa2222' },
};

function laidOut(nodes, edges) {
  const g = new dagre.graphlib.Graph();
  g.setGraph({ rankdir: 'LR', nodesep: 40, ranksep: 80 });
  g.setDefaultEdgeLabel(() => ({}));
  nodes.forEach((n) => g.setNode(n.id, { width: NODE_W, height: NODE_H }));
  edges.forEach((e) => g.setEdge(e.source, e.target));
  dagre.layout(g);
  return nodes.map((n) => {
    const { x, y } = g.node(n.id);
    return { ...n, position: { x: x - NODE_W / 2, y: y - NODE_H / 2 } };
  });
}

function buildGraph(model, layers) {
  if (!model) return { nodes: [], edges: [] };

  const satisfyingComponents = new Map();
  const allocatedReqs = new Map();
  for (const c of model.components) {
    for (const s of c.Satisfies || []) {
      if (!satisfyingComponents.has(s)) satisfyingComponents.set(s, []);
      satisfyingComponents.get(s).push(c.ID);
    }
  }
  for (const r of model.requirements) {
    if (r.AllocatedTo) {
      if (!allocatedReqs.has(r.AllocatedTo)) allocatedReqs.set(r.AllocatedTo, []);
      allocatedReqs.get(r.AllocatedTo).push(r.ID);
    }
  }

  const nodes = [];
  const edges = [];

  for (const c of model.components) {
    nodes.push({
      id: c.ID,
      type: 'default',
      data: { label: `${c.ID}\n${c.Name}`, kind: 'component', payload: c },
      style: nodeStyle('component'),
      position: { x: 0, y: 0 },
    });
  }

  for (const r of model.requirements) {
    const isSystem = !r.DerivedFrom;
    const isUnallocated = isSystem && !(satisfyingComponents.get(r.ID) || []).length;
    if (isSystem) {
      nodes.push({
        id: r.ID,
        type: 'default',
        data: { label: `${r.ID}\n${r.Title}`, kind: 'sysreq', payload: r },
        style: nodeStyle(isUnallocated ? 'gap' : 'sysreq'),
        position: { x: 0, y: 0 },
      });
    } else if (layers.derivedReqs) {
      nodes.push({
        id: r.ID,
        type: 'default',
        data: { label: `${r.ID}\n${r.Title}`, kind: 'derivedReq', payload: r },
        style: nodeStyle('derivedReq'),
        position: { x: 0, y: 0 },
      });
    }
  }

  if (layers.testSpecs) {
    for (const t of model.test_specs) {
      nodes.push({
        id: t.ID,
        type: 'default',
        data: { label: `${t.ID}\n${t.Title}`, kind: 'testspec', payload: t },
        style: nodeStyle('testspec'),
        position: { x: 0, y: 0 },
      });
    }
  }

  if (layers.interfaces) {
    for (const i of model.interfaces) {
      nodes.push({
        id: `iface:${i.Name}`,
        type: 'default',
        data: { label: i.Name, kind: 'interface', payload: i },
        style: nodeStyle('interface'),
        position: { x: 0, y: 0 },
      });
    }
  }

  for (const c of model.components) {
    for (const s of c.Satisfies || []) {
      edges.push({ id: `${c.ID}->sat->${s}`, source: c.ID, target: s, label: 'satisfies' });
    }
    for (const d of c.DependsOn || []) {
      edges.push({
        id: `${c.ID}->dep->${d}`,
        source: c.ID,
        target: d,
        label: 'depends_on',
        style: { strokeDasharray: '4 4' },
      });
    }
  }

  if (layers.derivedReqs) {
    for (const r of model.requirements) {
      if (r.DerivedFrom && nodes.some((n) => n.id === r.DerivedFrom)) {
        edges.push({
          id: `${r.ID}->df->${r.DerivedFrom}`,
          source: r.ID,
          target: r.DerivedFrom,
          label: 'derived_from',
          style: { strokeDasharray: '4 4' },
        });
      }
      if (r.AllocatedTo && nodes.some((n) => n.id === r.AllocatedTo)) {
        edges.push({
          id: `${r.ID}->alloc->${r.AllocatedTo}`,
          source: r.ID,
          target: r.AllocatedTo,
          label: 'allocated_to',
        });
      }
    }
  }

  if (layers.testSpecs) {
    for (const t of model.test_specs) {
      for (const v of t.Verifies || []) {
        if (nodes.some((n) => n.id === v)) {
          edges.push({ id: `${t.ID}->ver->${v}`, source: t.ID, target: v, label: 'verifies' });
        }
      }
    }
  }

  return { nodes: laidOut(nodes, edges), edges };
}

function nodeStyle(kind) {
  const palette = TYPE_STYLE[kind];
  return {
    background: palette.background,
    border: `1.5px solid ${palette.border}`,
    borderRadius: 6,
    fontSize: 12,
    padding: 8,
    width: NODE_W,
    textAlign: 'center',
    whiteSpace: 'pre-line',
  };
}

function DetailPanel({ selected }) {
  if (!selected) {
    return (
      <div style={panelStyle}>
        <h3 style={{ margin: 0, fontSize: 14 }}>Click a node</h3>
        <p style={{ fontSize: 12, color: '#666' }}>
          Node details, source file, and traceability will appear here.
        </p>
      </div>
    );
  }
  const { data } = selected;
  const p = data.payload;
  return (
    <div style={panelStyle}>
      <h3 style={{ margin: 0, fontSize: 14 }}>
        {p.ID || p.Name} <span style={{ color: '#888', fontWeight: 400 }}>({data.kind})</span>
      </h3>
      <p style={{ fontSize: 12, color: '#222', marginTop: 6 }}>{p.Title || p.Name || ''}</p>
      {p.Statement && <div style={fieldBlock}><b>Statement:</b><br/>{p.Statement}</div>}
      {p.Responsibility && <div style={fieldBlock}><b>Responsibility:</b><br/>{p.Responsibility}</div>}
      {p.Priority && <div style={fieldRow}><b>Priority:</b> {p.Priority}</div>}
      {p.Pattern && <div style={fieldRow}><b>Pattern:</b> {p.Pattern}</div>}
      {p.DerivedFrom && <div style={fieldRow}><b>Derived from:</b> {p.DerivedFrom}</div>}
      {p.AllocatedTo && <div style={fieldRow}><b>Allocated to:</b> {p.AllocatedTo}</div>}
      {p.Satisfies && p.Satisfies.length > 0 && <div style={fieldRow}><b>Satisfies:</b> {p.Satisfies.join(', ')}</div>}
      {p.DependsOn && p.DependsOn.length > 0 && <div style={fieldRow}><b>Depends on:</b> {p.DependsOn.join(', ')}</div>}
      {p.Verifies && p.Verifies.length > 0 && <div style={fieldRow}><b>Verifies:</b> {p.Verifies.join(', ')}</div>}
      {p.Given && <div style={fieldBlock}><b>Given:</b><br/>{p.Given}</div>}
      {p.Expect && <div style={fieldBlock}><b>Expect:</b><br/>{p.Expect}</div>}
      {p.SourceFile && (
        <div style={{ ...fieldRow, marginTop: 12, color: '#666' }}>
          <b>Source:</b> {p.SourceFile}:{p.LineNumber}
        </div>
      )}
    </div>
  );
}

const panelStyle = {
  position: 'absolute',
  top: 12,
  right: 12,
  width: 340,
  padding: 14,
  background: '#ffffff',
  border: '1px solid #ddd',
  borderRadius: 8,
  boxShadow: '0 2px 8px rgba(0,0,0,0.05)',
  fontSize: 12,
  zIndex: 10,
  maxHeight: 'calc(100vh - 24px)',
  overflowY: 'auto',
};

const fieldRow = { marginTop: 6, fontSize: 12 };
const fieldBlock = { marginTop: 10, fontSize: 12, background: '#fafafa', padding: 8, borderRadius: 4 };

function LayerControls({ layers, setLayers }) {
  const Toggle = ({ k, label }) => (
    <label style={{ display: 'flex', alignItems: 'center', gap: 6, fontSize: 12 }}>
      <input
        type="checkbox"
        checked={layers[k]}
        onChange={() => setLayers({ ...layers, [k]: !layers[k] })}
      />
      {label}
    </label>
  );
  return (
    <div style={{ display: 'flex', gap: 16, background: '#fff', padding: '8px 12px', border: '1px solid #ddd', borderRadius: 6 }}>
      <Toggle k="derivedReqs" label="Derived reqs" />
      <Toggle k="testSpecs" label="Test specs" />
      <Toggle k="interfaces" label="Interfaces" />
    </div>
  );
}

export default function App() {
  const [model, setModel] = useState(null);
  const [error, setError] = useState(null);
  const [layers, setLayers] = useState({ derivedReqs: false, testSpecs: false, interfaces: false });
  const [selected, setSelected] = useState(null);

  useEffect(() => {
    fetch('/api/model')
      .then((r) => {
        if (!r.ok) throw new Error(`HTTP ${r.status}`);
        return r.json();
      })
      .then(setModel)
      .catch((e) => setError(String(e)));
  }, []);

  const graph = useMemo(() => buildGraph(model, layers), [model, layers]);
  const [nodes, setNodes, onNodesChange] = useNodesState(graph.nodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(graph.edges);

  useEffect(() => {
    setNodes(graph.nodes);
    setEdges(graph.edges);
  }, [graph, setNodes, setEdges]);

  const onNodeClick = useCallback((_, node) => setSelected(node), []);

  if (error) return <div style={{ padding: 20, color: '#a22' }}>Failed to load model: {error}</div>;
  if (!model) return <div style={{ padding: 20 }}>Loading model…</div>;

  return (
    <div style={{ width: '100vw', height: '100vh' }}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onNodeClick={onNodeClick}
        onPaneClick={() => setSelected(null)}
        fitView
      >
        <Background />
        <Controls />
        <MiniMap pannable zoomable />
        <Panel position="top-left">
          <LayerControls layers={layers} setLayers={setLayers} />
        </Panel>
      </ReactFlow>
      <DetailPanel selected={selected} />
    </div>
  );
}
