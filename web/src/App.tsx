import React, { useState, useEffect, useRef, useCallback } from 'react';
import { Template, TemplateSummary } from './lib/types';
import {
  listTemplates,
  getTemplate,
  saveTemplate,
  deleteTemplate,
  renderLivePreview,
} from './lib/api';
import { CanvasPreview } from './components/CanvasPreview';
import { TemplateSelector } from './components/TemplateSelector';
import { ElementInspector } from './components/ElementInspector';
import { ParamTester } from './components/ParamTester';
import { CodeExporter } from './components/CodeExporter';
import { SocialFeedSimulator } from './components/SocialFeedSimulator';
import { TemplateImportModal } from './components/TemplateImportModal';
import {
  Cpu,
  Save,
  ShieldCheck,
  Undo2,
  Redo2,
  Upload,
  Download,
  Eye,
  LayoutTemplate,
  Layers,
  Sliders,
} from 'lucide-react';

export const App: React.FC = () => {
  const [templateSummaries, setTemplateSummaries] = useState<TemplateSummary[]>([]);
  const [activeTemplateId, setActiveTemplateId] = useState<string>('dark_tech');
  const [activeTemplate, setActiveTemplate] = useState<Template | null>(null);
  const [selectedElementIndex, setSelectedElementIndex] = useState<number | null>(null);

  // Undo / Redo history
  const [history, setHistory] = useState<Template[]>([]);
  const [historyIndex, setHistoryIndex] = useState<number>(-1);

  // View Mode: 'canvas' or 'social-feed'
  const [viewMode, setViewMode] = useState<'canvas' | 'social-feed'>('canvas');

  // Mobile Navigation Tab: 'templates' | 'canvas' | 'inspector' | 'feeds'
  const [mobileTab, setMobileTab] = useState<'templates' | 'canvas' | 'inspector' | 'feeds'>('canvas');

  // Modal
  const [isImportModalOpen, setIsImportModalOpen] = useState(false);

  const [params, setParams] = useState<Record<string, string>>({});
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [renderTime, setRenderTime] = useState<string>('');
  const [isSaving, setIsSaving] = useState<boolean>(false);
  const [saveStatus, setSaveStatus] = useState<string | null>(null);

  const debounceTimerRef = useRef<number | null>(null);

  // 1. Fetch template summaries
  const loadTemplatesList = useCallback(async () => {
    try {
      const list = await listTemplates();
      setTemplateSummaries(list);
      if (list.length > 0 && !list.some((t) => t.id === activeTemplateId)) {
        setActiveTemplateId(list[0].id);
      }
    } catch (e) {
      console.error('Failed to list templates', e);
    }
  }, [activeTemplateId]);

  useEffect(() => {
    loadTemplatesList();
  }, [loadTemplatesList]);

  // 2. Fetch active template schema
  useEffect(() => {
    if (!activeTemplateId) return;

    let cancelled = false;
    getTemplate(activeTemplateId)
      .then((tmpl) => {
        if (cancelled) return;
        setActiveTemplate(tmpl);
        setHistory([tmpl]);
        setHistoryIndex(0);
        setSelectedElementIndex(null);

        const initialParams: Record<string, string> = {};
        for (const el of tmpl.elements) {
          if (el.field && el.default) {
            initialParams[el.field] = el.default;
          }
        }
        setParams(initialParams);
      })
      .catch((err) => console.error(err));

    return () => {
      cancelled = true;
    };
  }, [activeTemplateId]);

  // 3. Update active template with history push
  const handleUpdateTemplate = useCallback((newTemplate: Template) => {
    setActiveTemplate(newTemplate);
    setHistory((prev) => {
      const updated = prev.slice(0, historyIndex + 1);
      return [...updated, newTemplate];
    });
    setHistoryIndex((prev) => prev + 1);
  }, [historyIndex]);

  // Undo / Redo
  const handleUndo = useCallback(() => {
    if (historyIndex > 0) {
      const targetIndex = historyIndex - 1;
      setHistoryIndex(targetIndex);
      setActiveTemplate(history[targetIndex]);
    }
  }, [historyIndex, history]);

  const handleRedo = useCallback(() => {
    if (historyIndex < history.length - 1) {
      const targetIndex = historyIndex + 1;
      setHistoryIndex(targetIndex);
      setActiveTemplate(history[targetIndex]);
    }
  }, [historyIndex, history]);

  // Keyboard shortcuts (Ctrl+Z / Ctrl+Y)
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.ctrlKey || e.metaKey) && e.key === 'z') {
        if (e.shiftKey) {
          handleRedo();
        } else {
          handleUndo();
        }
      } else if ((e.ctrlKey || e.metaKey) && e.key === 'y') {
        handleRedo();
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [handleUndo, handleRedo]);

  // 4. Render card preview on template or params change
  const triggerRender = useCallback(
    (tmpl: Template, currentParams: Record<string, string>) => {
      if (debounceTimerRef.current) {
        clearTimeout(debounceTimerRef.current);
      }

      debounceTimerRef.current = window.setTimeout(async () => {
        setLoading(true);
        const startTime = performance.now();
        try {
          const url = await renderLivePreview(tmpl, currentParams);
          setPreviewUrl(url);
          const duration = (performance.now() - startTime).toFixed(1);
          setRenderTime(`${duration}ms`);
        } catch (e) {
          console.error('Failed to render live preview', e);
        } finally {
          setLoading(false);
        }
      }, 150);
    },
    []
  );

  useEffect(() => {
    if (!activeTemplate) return;
    triggerRender(activeTemplate, params);
  }, [activeTemplate, params, triggerRender]);

  const handleSaveCurrentTemplate = async () => {
    if (!activeTemplate) return;
    setIsSaving(true);
    try {
      await saveTemplate(activeTemplate);
      setSaveStatus('Template saved successfully!');
      setTimeout(() => setSaveStatus(null), 2500);
      loadTemplatesList();
    } catch (e: any) {
      alert(e.message || 'Failed to save template');
    } finally {
      setIsSaving(false);
    }
  };

  const handleExportJSON = () => {
    if (!activeTemplate) return;
    const blob = new Blob([JSON.stringify(activeTemplate, null, 2)], {
      type: 'application/json',
    });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${activeTemplate.id}.json`;
    a.click();
    URL.revokeObjectURL(url);
  };

  const handleDeleteCustomTemplate = async (id: string) => {
    try {
      await deleteTemplate(id);
      loadTemplatesList();
      if (activeTemplateId === id) {
        setActiveTemplateId('dark_tech');
      }
    } catch (e: any) {
      alert(e.message || 'Failed to delete template');
    }
  };

  const handleParamChange = (field: string, value: string) => {
    setParams((prev) => ({
      ...prev,
      [field]: value,
    }));
  };

  return (
    <div className="flex flex-col h-screen w-screen overflow-hidden bg-[#070a10] text-slate-100 font-sans">
      {/* Top Navbar */}
      <header className="h-14 border-b border-[#1a2333] bg-[#0c121e] px-5 flex items-center justify-between flex-shrink-0 z-20">
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-lg bg-sky-500/15 border border-sky-500/30 flex items-center justify-center font-mono font-bold text-sky-400">
              OG
            </div>
            <div>
              <div className="flex items-center gap-2">
                <span className="font-bold text-sm tracking-tight text-white">ogcraft</span>
                <span className="text-[10px] px-1.5 py-0.2 rounded bg-sky-500/20 text-sky-400 font-mono font-semibold">
                  v1.0
                </span>
              </div>
              <p className="text-[10px] text-slate-400 -mt-0.5">Real-Time OpenGraph Studio</p>
            </div>
          </div>

          <div className="h-5 w-px bg-slate-800 mx-1 hidden sm:block" />

          {/* Undo / Redo */}
          <div className="flex items-center gap-1 bg-slate-900 border border-slate-800 rounded-lg p-0.5">
            <button
              onClick={handleUndo}
              disabled={historyIndex <= 0}
              className="p-1.5 hover:bg-slate-800 text-slate-400 hover:text-slate-200 disabled:opacity-30 rounded transition"
              title="Undo (Ctrl+Z)"
            >
              <Undo2 size={13} />
            </button>
            <button
              onClick={handleRedo}
              disabled={historyIndex >= history.length - 1}
              className="p-1.5 hover:bg-slate-800 text-slate-400 hover:text-slate-200 disabled:opacity-30 rounded transition"
              title="Redo (Ctrl+Y)"
            >
              <Redo2 size={13} />
            </button>
          </div>

          <div className="hidden lg:flex items-center gap-4 text-xs font-mono text-slate-400 ml-2">
            <span className="flex items-center gap-1.5">
              <Cpu size={13} className="text-emerald-400" />
              SIMD Go 2D (WebP, PNG, JPEG)
            </span>
            <span className="flex items-center gap-1.5">
              <ShieldCheck size={13} className="text-sky-400" />
              SSRF & HMAC Protected
            </span>
          </div>
        </div>

        {/* View Mode Switcher */}
        <div className="flex items-center bg-[#070a10] border border-slate-800 rounded-lg p-0.5">
          <button
            onClick={() => setViewMode('canvas')}
            className={`flex items-center gap-1.5 px-3 py-1 text-xs font-medium rounded-md transition ${
              viewMode === 'canvas'
                ? 'bg-sky-500/20 text-sky-300 font-semibold shadow'
                : 'text-slate-400 hover:text-slate-200'
            }`}
          >
            <LayoutTemplate size={13} />
            Canvas Editor
          </button>
          <button
            onClick={() => setViewMode('social-feed')}
            className={`flex items-center gap-1.5 px-3 py-1 text-xs font-medium rounded-md transition ${
              viewMode === 'social-feed'
                ? 'bg-sky-500/20 text-sky-300 font-semibold shadow'
                : 'text-slate-400 hover:text-slate-200'
            }`}
          >
            <Eye size={13} />
            Social Feeds
          </button>
        </div>

        {/* Right Action Buttons */}
        <div className="flex items-center gap-2">
          {saveStatus && (
            <span className="text-xs text-emerald-400 font-mono transition">
              ✓ {saveStatus}
            </span>
          )}

          <button
            onClick={() => setIsImportModalOpen(true)}
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-slate-900 hover:bg-slate-800 text-slate-300 text-xs font-medium border border-slate-800 transition"
            title="Import JSON"
          >
            <Upload size={13} />
            <span className="hidden sm:inline">Import</span>
          </button>

          <button
            onClick={handleExportJSON}
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-slate-900 hover:bg-slate-800 text-slate-300 text-xs font-medium border border-slate-800 transition"
            title="Export JSON"
          >
            <Download size={13} />
            <span className="hidden sm:inline">Export</span>
          </button>

          <button
            onClick={handleSaveCurrentTemplate}
            disabled={isSaving || !activeTemplate}
            className="flex items-center gap-1.5 px-3.5 py-1.5 rounded-lg bg-sky-500 hover:bg-sky-400 text-slate-950 text-xs font-bold transition"
          >
            <Save size={13} />
            {isSaving ? 'Saving...' : 'Save Template'}
          </button>
        </div>
      </header>

      {/* Main Studio Body */}
      <div className="flex-1 flex overflow-hidden relative">
        {/* Left: Template Selector */}
        <div className={`${mobileTab === 'templates' ? 'flex' : 'hidden'} lg:flex h-full w-full lg:w-72 flex-shrink-0`}>
          <TemplateSelector
            templates={templateSummaries}
            activeId={activeTemplateId}
            onSelect={(id) => {
              setActiveTemplateId(id);
              if (window.innerWidth < 1024) {
                setMobileTab('canvas');
              }
            }}
            onCreateNew={() => {
              const newId = `custom_${Date.now().toString().slice(-4)}`;
              const tmpl: Template = {
                id: newId,
                name: 'New Custom Template',
                width: 1200,
                height: 630,
                background: {
                  type: 'gradient',
                  direction: 'to-bottom-right',
                  pattern: 'dots',
                  stops: [
                    { offset: 0, color: '#090d16' },
                    { offset: 1, color: '#161e2e' },
                  ],
                },
                elements: [
                  {
                    type: 'badge',
                    field: 'tag',
                    default: 'CUSTOM',
                    font: 'JetBrainsMono-Medium',
                    font_size: 18,
                    color: '#ffffff',
                    background_color: '#2563eb',
                    x: 80,
                    y: 80,
                    padding_x: 14,
                    padding_y: 6,
                    border_radius: 6,
                  },
                  {
                    type: 'text',
                    field: 'title',
                    default: 'Your Custom Headline Here',
                    font: 'Inter-Bold',
                    font_size: 54,
                    color: '#f8fafc',
                    x: 80,
                    y: 160,
                    max_width: 1040,
                    line_height: 1.25,
                    max_lines: 3,
                    fit_text: true,
                  },
                ],
              };
              handleUpdateTemplate(tmpl);
              setActiveTemplateId(newId);
              if (window.innerWidth < 1024) {
                setMobileTab('canvas');
              }
            }}
            onDelete={handleDeleteCustomTemplate}
          />
        </div>

        {/* Center Viewport */}
        {activeTemplate ? (
          <div
            className={`${
              mobileTab === 'canvas' || mobileTab === 'feeds' ? 'flex' : 'hidden'
            } lg:flex flex-1 flex-col h-full overflow-hidden`}
          >
            {viewMode === 'canvas' && mobileTab !== 'feeds' ? (
              <div className="flex-1 flex flex-col h-full overflow-hidden">
                <CanvasPreview
                  template={activeTemplate}
                  previewUrl={previewUrl}
                  loading={loading}
                  renderTime={renderTime}
                  selectedElementIndex={selectedElementIndex}
                  onSelectElement={setSelectedElementIndex}
                  onUpdateElement={(idx, updated) => {
                    const nextElements = [...activeTemplate.elements];
                    nextElements[idx] = updated;
                    handleUpdateTemplate({
                      ...activeTemplate,
                      elements: nextElements,
                    });
                  }}
                  onRefresh={() => triggerRender(activeTemplate, params)}
                />

                {/* Bottom Controls: Simulator & Exporter */}
                <div className="flex-shrink-0 flex flex-col max-h-72 overflow-y-auto">
                  <ParamTester
                    template={activeTemplate}
                    params={params}
                    onChangeParam={handleParamChange}
                    onResetParams={() => {
                      const initial: Record<string, string> = {};
                      for (const el of activeTemplate.elements) {
                        if (el.field && el.default) {
                          initial[el.field] = el.default;
                        }
                      }
                      setParams(initial);
                    }}
                  />
                  <CodeExporter template={activeTemplate} params={params} />
                </div>
              </div>
            ) : (
              <SocialFeedSimulator
                template={activeTemplate}
                previewUrl={previewUrl}
                params={params}
              />
            )}
          </div>
        ) : (
          <div className="flex-1 flex items-center justify-center text-slate-500 font-mono text-sm">
            Loading template...
          </div>
        )}

        {/* Right Inspector */}
        {activeTemplate && (
          <div
            className={`${
              mobileTab === 'inspector' ? 'flex' : 'hidden'
            } ${viewMode === 'canvas' ? 'lg:flex' : 'lg:hidden'} h-full w-full lg:w-80 flex-shrink-0`}
          >
            <ElementInspector
              template={activeTemplate}
              selectedElementIndex={selectedElementIndex}
              onUpdateTemplate={handleUpdateTemplate}
              onSelectElement={setSelectedElementIndex}
            />
          </div>
        )}
      </div>

      {/* Mobile Bottom Navigation Tab Bar */}
      <nav className="lg:hidden h-14 bg-[#0c121e] border-t border-[#1a2333] px-2 flex items-center justify-around flex-shrink-0 z-30">
        <button
          onClick={() => setMobileTab('templates')}
          className={`flex flex-col items-center justify-center gap-1 flex-1 py-1 rounded-lg text-[10px] font-medium transition ${
            mobileTab === 'templates'
              ? 'text-sky-400 bg-sky-500/10 font-bold'
              : 'text-slate-400 hover:text-slate-200'
          }`}
        >
          <LayoutTemplate size={16} />
          <span>Templates</span>
        </button>
        <button
          onClick={() => {
            setMobileTab('canvas');
            setViewMode('canvas');
          }}
          className={`flex flex-col items-center justify-center gap-1 flex-1 py-1 rounded-lg text-[10px] font-medium transition ${
            mobileTab === 'canvas' && viewMode === 'canvas'
              ? 'text-sky-400 bg-sky-500/10 font-bold'
              : 'text-slate-400 hover:text-slate-200'
          }`}
        >
          <Layers size={16} />
          <span>Canvas</span>
        </button>
        <button
          onClick={() => {
            setMobileTab('inspector');
            setViewMode('canvas');
          }}
          className={`flex flex-col items-center justify-center gap-1 flex-1 py-1 rounded-lg text-[10px] font-medium transition ${
            mobileTab === 'inspector'
              ? 'text-sky-400 bg-sky-500/10 font-bold'
              : 'text-slate-400 hover:text-slate-200'
          }`}
        >
          <Sliders size={16} />
          <span>Inspector</span>
        </button>
        <button
          onClick={() => {
            setMobileTab('feeds');
            setViewMode('social-feed');
          }}
          className={`flex flex-col items-center justify-center gap-1 flex-1 py-1 rounded-lg text-[10px] font-medium transition ${
            mobileTab === 'feeds' || viewMode === 'social-feed'
              ? 'text-sky-400 bg-sky-500/10 font-bold'
              : 'text-slate-400 hover:text-slate-200'
          }`}
        >
          <Eye size={16} />
          <span>Feeds</span>
        </button>
      </nav>

      {/* Import Modal */}
      <TemplateImportModal
        isOpen={isImportModalOpen}
        onClose={() => setIsImportModalOpen(false)}
        onImport={(imported) => {
          handleUpdateTemplate(imported);
          setActiveTemplateId(imported.id);
        }}
      />
    </div>
  );
};
