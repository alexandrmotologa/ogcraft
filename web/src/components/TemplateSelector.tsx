import React from 'react';
import { TemplateSummary } from '../lib/types';
import { Layout, Plus, Trash2, Sparkles } from 'lucide-react';

interface TemplateSelectorProps {
  templates: TemplateSummary[];
  activeId: string;
  onSelect: (id: string) => void;
  onCreateNew: () => void;
  onDelete: (id: string) => void;
}

export const TemplateSelector: React.FC<TemplateSelectorProps> = ({
  templates,
  activeId,
  onSelect,
  onCreateNew,
  onDelete,
}) => {
  return (
    <div className="flex flex-col h-full bg-[#0c121e] border-r border-[#1a2333] w-72 flex-shrink-0">
      {/* Header */}
      <div className="flex items-center justify-between p-4 border-b border-[#1a2333]">
        <div className="flex items-center gap-2">
          <Layout size={16} className="text-sky-400" />
          <span className="text-xs font-bold uppercase tracking-wider text-slate-300">
            Templates ({templates.length})
          </span>
        </div>
        <button
          onClick={onCreateNew}
          className="p-1.5 hover:bg-slate-800 text-sky-400 hover:text-sky-300 rounded transition"
          title="New custom template"
        >
          <Plus size={16} />
        </button>
      </div>

      {/* List */}
      <div className="flex-1 overflow-y-auto p-2 space-y-1">
        {templates.map((t) => {
          const isActive = t.id === activeId;
          return (
            <div
              key={t.id}
              onClick={() => onSelect(t.id)}
              className={`group flex items-start justify-between p-2.5 rounded-lg cursor-pointer transition ${
                isActive
                  ? 'bg-sky-500/15 border border-sky-500/40 text-white'
                  : 'hover:bg-slate-800/60 text-slate-400 hover:text-slate-200 border border-transparent'
              }`}
            >
              <div className="flex flex-col min-w-0 pr-2">
                <div className="flex items-center gap-2">
                  <span className={`text-xs font-semibold truncate ${isActive ? 'text-sky-300' : 'text-slate-200'}`}>
                    {t.name}
                  </span>
                  {t.builtin && (
                    <span className="text-[9px] px-1.5 py-0.2 rounded bg-slate-800 text-slate-400 font-mono">
                      PRESET
                    </span>
                  )}
                </div>
                <span className="text-[11px] text-slate-500 truncate mt-0.5 font-mono">
                  {t.id} • {t.elements_count} items
                </span>
              </div>

              {!t.builtin && (
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    if (confirm(`Delete custom template ${t.name}?`)) {
                      onDelete(t.id);
                    }
                  }}
                  className="opacity-0 group-hover:opacity-100 p-1 hover:bg-rose-500/20 text-slate-500 hover:text-rose-400 rounded transition"
                  title="Delete custom template"
                >
                  <Trash2 size={13} />
                </button>
              )}
            </div>
          );
        })}
      </div>

      {/* Builtin info notice */}
      <div className="p-3 border-t border-[#1a2333] bg-[#090d16]/60">
        <div className="flex items-center gap-1.5 text-[11px] text-slate-400 font-mono">
          <Sparkles size={12} className="text-amber-400 flex-shrink-0" />
          <span>SIMD vector rendering &lt; 15ms</span>
        </div>
      </div>
    </div>
  );
};
