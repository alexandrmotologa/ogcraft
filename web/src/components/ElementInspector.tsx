import React from 'react';
import { ElementSpec, Template } from '../lib/types';
import { Sliders, Trash2, Copy, Plus, Type, ShieldAlert, Image, Sparkles, Minus, Square } from 'lucide-react';

interface ElementInspectorProps {
  template: Template;
  selectedElementIndex: number | null;
  onUpdateTemplate: (updated: Template) => void;
  onSelectElement: (index: number | null) => void;
}

export const ElementInspector: React.FC<ElementInspectorProps> = ({
  template,
  selectedElementIndex,
  onUpdateTemplate,
  onSelectElement,
}) => {
  const selectedElement =
    selectedElementIndex !== null ? template.elements[selectedElementIndex] : null;

  const updateSelectedElement = (partial: Partial<ElementSpec>) => {
    if (selectedElementIndex === null || !selectedElement) return;
    const updatedElements = [...template.elements];
    updatedElements[selectedElementIndex] = {
      ...selectedElement,
      ...partial,
    };
    onUpdateTemplate({
      ...template,
      elements: updatedElements,
    });
  };

  const deleteSelectedElement = () => {
    if (selectedElementIndex === null) return;
    const updatedElements = template.elements.filter((_, i) => i !== selectedElementIndex);
    onSelectElement(null);
    onUpdateTemplate({
      ...template,
      elements: updatedElements,
    });
  };

  const duplicateSelectedElement = () => {
    if (selectedElementIndex === null || !selectedElement) return;
    const copy: ElementSpec = {
      ...selectedElement,
      x: selectedElement.x + 20,
      y: selectedElement.y + 20,
    };
    const updatedElements = [...template.elements, copy];
    onUpdateTemplate({
      ...template,
      elements: updatedElements,
    });
    onSelectElement(updatedElements.length - 1);
  };

  const addElement = (type: ElementSpec['type']) => {
    let newEl: ElementSpec;
    switch (type) {
      case 'text':
        newEl = {
          type: 'text',
          field: 'title',
          default: 'New Headline Text',
          font: 'Inter-Bold',
          font_size: 48,
          color: '#f8fafc',
          x: 80,
          y: 200,
          max_width: 1040,
          line_height: 1.25,
          max_lines: 3,
          fit_text: true,
        };
        break;
      case 'badge':
        newEl = {
          type: 'badge',
          field: 'tag',
          default: 'NEW TAG',
          font: 'JetBrainsMono-Medium',
          font_size: 20,
          color: '#ffffff',
          background_color: '#3b82f6',
          x: 80,
          y: 90,
          padding_x: 14,
          padding_y: 7,
          border_radius: 6,
        };
        break;
      case 'avatar':
        newEl = {
          type: 'avatar',
          field: 'avatar',
          default: 'https://avatars.githubusercontent.com/u/583231',
          x: 80,
          y: 480,
          size: 64,
          shape: 'circle',
          border_color: '#334155',
          border_width: 2,
        };
        break;
      case 'icon':
        newEl = {
          type: 'icon',
          icon: 'star',
          size: 32,
          color: '#facc15',
          x: 80,
          y: 490,
        };
        break;
      case 'divider':
        newEl = {
          type: 'divider',
          x: 80,
          y: 440,
          width: 1040,
          height: 1,
          color: '#1e293b',
        };
        break;
      case 'container':
        newEl = {
          type: 'container',
          x: 60,
          y: 60,
          width: 1080,
          height: 510,
          border_radius: 16,
          background_color: '#0f172a88',
          border_color: '#1e293b',
          border_width: 1,
        };
        break;
    }

    const updatedElements = [...template.elements, newEl];
    onUpdateTemplate({
      ...template,
      elements: updatedElements,
    });
    onSelectElement(updatedElements.length - 1);
  };

  return (
    <div className="flex flex-col h-full bg-[#0c121e] border-l border-[#1a2333] w-80 flex-shrink-0 overflow-y-auto">
      {/* Top Header */}
      <div className="flex items-center justify-between p-4 border-b border-[#1a2333]">
        <div className="flex items-center gap-2">
          <Sliders size={16} className="text-sky-400" />
          <span className="text-xs font-bold uppercase tracking-wider text-slate-300">
            {selectedElement ? `${selectedElement.type} Inspector` : 'Canvas Style'}
          </span>
        </div>

        {selectedElement && (
          <div className="flex items-center gap-1">
            <button
              onClick={duplicateSelectedElement}
              className="p-1 hover:bg-slate-800 text-slate-400 hover:text-slate-200 rounded transition"
              title="Duplicate element"
            >
              <Copy size={14} />
            </button>
            <button
              onClick={deleteSelectedElement}
              className="p-1 hover:bg-rose-500/20 text-slate-400 hover:text-rose-400 rounded transition"
              title="Delete element"
            >
              <Trash2 size={14} />
            </button>
          </div>
        )}
      </div>

      <div className="p-4 space-y-5 flex-1">
        {selectedElement ? (
          <>
            {/* Position */}
            <div className="space-y-2">
              <label className="text-[11px] font-bold text-slate-400 uppercase tracking-wider">
                Position (X, Y)
              </label>
              <div className="grid grid-cols-2 gap-2">
                <div>
                  <span className="text-[10px] text-slate-500 font-mono">X (px)</span>
                  <input
                    type="number"
                    value={selectedElement.x}
                    onChange={(e) => updateSelectedElement({ x: Number(e.target.value) })}
                    className="w-full bg-[#070a10] border border-slate-800 rounded px-2.5 py-1.5 text-xs text-slate-200 font-mono"
                  />
                </div>
                <div>
                  <span className="text-[10px] text-slate-500 font-mono">Y (px)</span>
                  <input
                    type="number"
                    value={selectedElement.y}
                    onChange={(e) => updateSelectedElement({ y: Number(e.target.value) })}
                    className="w-full bg-[#070a10] border border-slate-800 rounded px-2.5 py-1.5 text-xs text-slate-200 font-mono"
                  />
                </div>
              </div>
            </div>

            {/* Field & Default for text/badge/avatar */}
            {(selectedElement.type === 'text' ||
              selectedElement.type === 'badge' ||
              selectedElement.type === 'avatar') && (
              <div className="space-y-3 pt-2 border-t border-slate-800/60">
                <div>
                  <label className="text-[11px] font-bold text-slate-400 uppercase tracking-wider block mb-1">
                    Query Parameter Key
                  </label>
                  <input
                    type="text"
                    value={selectedElement.field || ''}
                    placeholder="e.g. title, author, tag"
                    onChange={(e) => updateSelectedElement({ field: e.target.value })}
                    className="w-full bg-[#070a10] border border-slate-800 rounded px-2.5 py-1.5 text-xs text-slate-200 font-mono"
                  />
                </div>
                <div>
                  <label className="text-[11px] font-bold text-slate-400 uppercase tracking-wider block mb-1">
                    Default Fallback Value
                  </label>
                  <input
                    type="text"
                    value={selectedElement.default || ''}
                    placeholder="Fallback text or URL"
                    onChange={(e) => updateSelectedElement({ default: e.target.value })}
                    className="w-full bg-[#070a10] border border-slate-800 rounded px-2.5 py-1.5 text-xs text-slate-200 font-mono"
                  />
                </div>
              </div>
            )}

            {/* Typography options */}
            {selectedElement.type === 'text' && (
              <div className="space-y-3 pt-2 border-t border-slate-800/60">
                <div>
                  <label className="text-[11px] font-bold text-slate-400 uppercase tracking-wider block mb-1">
                    Font Family
                  </label>
                  <select
                    value={selectedElement.font || 'Inter-Bold'}
                    onChange={(e) => updateSelectedElement({ font: e.target.value })}
                    className="w-full bg-[#070a10] border border-slate-800 rounded px-2.5 py-1.5 text-xs text-slate-200 font-mono"
                  >
                    <option value="Inter-Regular">Inter Regular</option>
                    <option value="Inter-SemiBold">Inter SemiBold</option>
                    <option value="Inter-Bold">Inter Bold</option>
                    <option value="JetBrainsMono-Medium">JetBrains Mono Medium</option>
                    <option value="JetBrainsMono-Bold">JetBrains Mono Bold</option>
                  </select>
                </div>

                <div className="grid grid-cols-2 gap-2">
                  <div>
                    <span className="text-[10px] text-slate-500 font-mono">Font Size</span>
                    <input
                      type="number"
                      value={selectedElement.font_size || 48}
                      onChange={(e) => updateSelectedElement({ font_size: Number(e.target.value) })}
                      className="w-full bg-[#070a10] border border-slate-800 rounded px-2.5 py-1.5 text-xs text-slate-200 font-mono"
                    />
                  </div>
                  <div>
                    <span className="text-[10px] text-slate-500 font-mono">Max Lines</span>
                    <input
                      type="number"
                      value={selectedElement.max_lines || 3}
                      onChange={(e) => updateSelectedElement({ max_lines: Number(e.target.value) })}
                      className="w-full bg-[#070a10] border border-slate-800 rounded px-2.5 py-1.5 text-xs text-slate-200 font-mono"
                    />
                  </div>
                </div>

                <div>
                  <label className="text-[11px] font-bold text-slate-400 uppercase tracking-wider block mb-1">
                    Text Color
                  </label>
                  <div className="flex items-center gap-2">
                    <input
                      type="color"
                      value={selectedElement.color?.startsWith('#') ? selectedElement.color.slice(0, 7) : '#ffffff'}
                      onChange={(e) => updateSelectedElement({ color: e.target.value })}
                      className="w-7 h-7 rounded border border-slate-700 bg-transparent cursor-pointer"
                    />
                    <input
                      type="text"
                      value={selectedElement.color || '#ffffff'}
                      onChange={(e) => updateSelectedElement({ color: e.target.value })}
                      className="flex-1 bg-[#070a10] border border-slate-800 rounded px-2.5 py-1.5 text-xs text-slate-200 font-mono"
                    />
                  </div>
                </div>

                <div className="flex items-center justify-between pt-1">
                  <span className="text-xs text-slate-300">FitText Auto-Downscale</span>
                  <input
                    type="checkbox"
                    checked={selectedElement.fit_text !== false}
                    onChange={(e) => updateSelectedElement({ fit_text: e.target.checked })}
                    className="rounded bg-slate-900 border-slate-700 text-sky-500"
                  />
                </div>
              </div>
            )}

            {/* Badge specific */}
            {selectedElement.type === 'badge' && (
              <div className="space-y-3 pt-2 border-t border-slate-800/60">
                <div>
                  <label className="text-[11px] font-bold text-slate-400 uppercase tracking-wider block mb-1">
                    Background Color
                  </label>
                  <div className="flex items-center gap-2">
                    <input
                      type="color"
                      value={selectedElement.background_color?.startsWith('#') ? selectedElement.background_color.slice(0, 7) : '#3b82f6'}
                      onChange={(e) => updateSelectedElement({ background_color: e.target.value })}
                      className="w-7 h-7 rounded border border-slate-700 bg-transparent cursor-pointer"
                    />
                    <input
                      type="text"
                      value={selectedElement.background_color || '#3b82f6'}
                      onChange={(e) => updateSelectedElement({ background_color: e.target.value })}
                      className="flex-1 bg-[#070a10] border border-slate-800 rounded px-2.5 py-1.5 text-xs text-slate-200 font-mono"
                    />
                  </div>
                </div>

                <div>
                  <label className="text-[11px] font-bold text-slate-400 uppercase tracking-wider block mb-1">
                    Icon (optional)
                  </label>
                  <select
                    value={selectedElement.icon || ''}
                    onChange={(e) => updateSelectedElement({ icon: e.target.value })}
                    className="w-full bg-[#070a10] border border-slate-800 rounded px-2.5 py-1.5 text-xs text-slate-200 font-mono"
                  >
                    <option value="">None</option>
                    <option value="github">GitHub</option>
                    <option value="x">X / Twitter</option>
                    <option value="clock">Clock</option>
                    <option value="calendar">Calendar</option>
                    <option value="star">Star</option>
                    <option value="tag">Tag</option>
                    <option value="bolt">Bolt</option>
                    <option value="check">Check</option>
                  </select>
                </div>
              </div>
            )}

            {/* Avatar specific */}
            {selectedElement.type === 'avatar' && (
              <div className="space-y-3 pt-2 border-t border-slate-800/60">
                <div className="grid grid-cols-2 gap-2">
                  <div>
                    <span className="text-[10px] text-slate-500 font-mono">Size (px)</span>
                    <input
                      type="number"
                      value={selectedElement.size || 64}
                      onChange={(e) => updateSelectedElement({ size: Number(e.target.value) })}
                      className="w-full bg-[#070a10] border border-slate-800 rounded px-2.5 py-1.5 text-xs text-slate-200 font-mono"
                    />
                  </div>
                  <div>
                    <span className="text-[10px] text-slate-500 font-mono">Shape</span>
                    <select
                      value={selectedElement.shape || 'circle'}
                      onChange={(e) => updateSelectedElement({ shape: e.target.value as any })}
                      className="w-full bg-[#070a10] border border-slate-800 rounded px-2.5 py-1.5 text-xs text-slate-200 font-mono"
                    >
                      <option value="circle">Circle</option>
                      <option value="rounded">Rounded</option>
                      <option value="square">Square</option>
                    </select>
                  </div>
                </div>
              </div>
            )}
          </>
        ) : (
          /* Template Background Controls */
          <div className="space-y-4">
            <div>
              <label className="text-[11px] font-bold text-slate-400 uppercase tracking-wider block mb-1">
                Background Type
              </label>
              <select
                value={template.background.type}
                onChange={(e) =>
                  onUpdateTemplate({
                    ...template,
                    background: {
                      ...template.background,
                      type: e.target.value as any,
                    },
                  })
                }
                className="w-full bg-[#070a10] border border-slate-800 rounded px-2.5 py-1.5 text-xs text-slate-200 font-mono"
              >
                <option value="gradient">Linear Gradient</option>
                <option value="radial">Radial Glow</option>
                <option value="solid">Solid Color</option>
              </select>
            </div>

            {template.background.type === 'solid' && (
              <div>
                <label className="text-[11px] font-bold text-slate-400 uppercase tracking-wider block mb-1">
                  Color
                </label>
                <div className="flex items-center gap-2">
                  <input
                    type="color"
                    value={template.background.color || '#090d16'}
                    onChange={(e) =>
                      onUpdateTemplate({
                        ...template,
                        background: { ...template.background, color: e.target.value },
                      })
                    }
                    className="w-7 h-7 rounded border border-slate-700 bg-transparent cursor-pointer"
                  />
                  <input
                    type="text"
                    value={template.background.color || '#090d16'}
                    onChange={(e) =>
                      onUpdateTemplate({
                        ...template,
                        background: { ...template.background, color: e.target.value },
                      })
                    }
                    className="flex-1 bg-[#070a10] border border-slate-800 rounded px-2.5 py-1.5 text-xs text-slate-200 font-mono"
                  />
                </div>
              </div>
            )}

            {/* Quick Add Elements Palette */}
            <div className="pt-4 border-t border-slate-800">
              <label className="text-[11px] font-bold text-slate-400 uppercase tracking-wider block mb-2">
                Add Visual Elements
              </label>
              <div className="grid grid-cols-2 gap-2">
                <button
                  onClick={() => addElement('text')}
                  className="flex items-center gap-2 p-2 rounded bg-slate-900 border border-slate-800 hover:border-slate-700 text-xs text-slate-200 transition"
                >
                  <Type size={13} className="text-sky-400" />
                  + Text
                </button>
                <button
                  onClick={() => addElement('badge')}
                  className="flex items-center gap-2 p-2 rounded bg-slate-900 border border-slate-800 hover:border-slate-700 text-xs text-slate-200 transition"
                >
                  <Sparkles size={13} className="text-emerald-400" />
                  + Badge
                </button>
                <button
                  onClick={() => addElement('avatar')}
                  className="flex items-center gap-2 p-2 rounded bg-slate-900 border border-slate-800 hover:border-slate-700 text-xs text-slate-200 transition"
                >
                  <Image size={13} className="text-purple-400" />
                  + Avatar
                </button>
                <button
                  onClick={() => addElement('divider')}
                  className="flex items-center gap-2 p-2 rounded bg-slate-900 border border-slate-800 hover:border-slate-700 text-xs text-slate-200 transition"
                >
                  <Minus size={13} className="text-amber-400" />
                  + Divider
                </button>
                <button
                  onClick={() => addElement('container')}
                  className="flex items-center gap-2 p-2 rounded bg-slate-900 border border-slate-800 hover:border-slate-700 text-xs text-slate-200 transition col-span-2"
                >
                  <Square size={13} className="text-blue-400" />
                  + Card Container Frame
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};
