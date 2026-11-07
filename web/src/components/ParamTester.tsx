import React from 'react';
import { Template } from '../lib/types';
import { Variable, Sparkles } from 'lucide-react';

interface ParamTesterProps {
  template: Template;
  params: Record<string, string>;
  onChangeParam: (field: string, value: string) => void;
  onResetParams: () => void;
}

export const ParamTester: React.FC<ParamTesterProps> = ({
  template,
  params,
  onChangeParam,
  onResetParams,
}) => {
  // Extract all unique fields from template elements
  const fields = Array.from(
    new Set(
      template.elements
        .map((el) => el.field)
        .filter((f): f is string => Boolean(f && f.trim() !== ''))
    )
  );

  return (
    <div className="bg-[#0c121e] border-t border-[#1a2333] p-4 flex flex-col gap-3">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Variable size={15} className="text-sky-400" />
          <span className="text-xs font-bold uppercase tracking-wider text-slate-300">
            Live URL Parameter Simulator
          </span>
        </div>
        <button
          onClick={onResetParams}
          className="text-[11px] text-slate-500 hover:text-slate-300 font-mono transition"
        >
          Reset to Defaults
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3">
        {fields.map((field) => (
          <div key={field} className="flex flex-col gap-1">
            <span className="text-[11px] font-mono text-slate-400 font-semibold">
              ?{field}=
            </span>
            <input
              type="text"
              value={params[field] ?? ''}
              placeholder={`Enter ${field}...`}
              onChange={(e) => onChangeParam(field, e.target.value)}
              className="bg-[#070a10] border border-slate-800 rounded px-2.5 py-1.5 text-xs text-slate-200 font-mono focus:border-sky-500/50 focus:outline-none transition"
            />
          </div>
        ))}
      </div>
    </div>
  );
};
