import * as React from "react";
import { CheckIcon } from "lucide-react";
import { cn } from "@/lib/utils";
import { Input } from "@/components/ui/input";

type Option = {
  label: string;
  value: string;
};

interface AutocompleteInputProps {
  value?: string; // text in the input
  onValueChange: (value: string) => void; // called on typing + select + create
  options: Option[];
  onCreateOption?: (label: string) => void;

  placeholder?: string;
  className?: string;
  disabled?: boolean;
}

type ListItem =
  | { type: "option"; option: Option }
  | { type: "create"; label: string };

export function AutocompleteInput({
  value,
  onValueChange,
  options,
  onCreateOption,
  placeholder = "Type to search...",
  className,
  disabled,
}: AutocompleteInputProps) {
  const [open, setOpen] = React.useState(false);
  const [highlightedIndex, setHighlightedIndex] = React.useState<number | null>(
    null
  );
  const containerRef = React.useRef<HTMLDivElement | null>(null);
  const listRef = React.useRef<HTMLUListElement | null>(null);

  const current = value ?? "";

  const filtered = React.useMemo(
    () =>
      options.filter((o) =>
        o.label.toLowerCase().includes(current.toLowerCase())
      ),
    [options, current]
  );

  const canCreate =
    !!onCreateOption &&
    current.trim().length > 0 &&
    !options.some(
      (o) => o.label.toLowerCase() === current.trim().toLowerCase()
    );

  const items: ListItem[] = React.useMemo(() => {
    const base: ListItem[] = filtered.map((option) => ({
      type: "option",
      option,
    }));
    if (canCreate) {
      base.push({ type: "create", label: current.trim() });
    }
    return base;
  }, [filtered, canCreate, current]);

  const handleSelect = (opt: Option) => {
    onValueChange(opt.label); // write label into the input
    setOpen(false);
    setHighlightedIndex(null);
  };

  const handleCreate = () => {
    const label = current.trim();
    if (!label || !onCreateOption) return;

    onCreateOption(label);
    onValueChange(label);
    setOpen(false);
    setHighlightedIndex(null);
  };

  // Close on outside click
  React.useEffect(() => {
    if (!open) return;
    const handler = (e: MouseEvent) => {
      if (!containerRef.current) return;
      if (!containerRef.current.contains(e.target as Node)) {
        setOpen(false);
        setHighlightedIndex(null);
      }
    };
    document.addEventListener("mousedown", handler);
    return () => document.removeEventListener("mousedown", handler);
  }, [open]);

  const moveHighlight = (direction: 1 | -1) => {
    if (!items.length) return;
    setOpen(true);
    setHighlightedIndex((prev) => {
      if (prev === null) return direction === 1 ? 0 : items.length - 1;
      let next = prev + direction;
      if (next < 0) next = items.length - 1;
      if (next >= items.length) next = 0;
      return next;
    });
  };

  // keep highlighted item in view
  React.useEffect(() => {
    if (highlightedIndex === null || !listRef.current) return;
    const el = listRef.current.children[highlightedIndex] as HTMLElement | null;
    if (el) {
      el.scrollIntoView({ block: "nearest" });
    }
  }, [highlightedIndex]);

  return (
    <div ref={containerRef} className={cn("relative", className)}>
      <Input
        disabled={disabled}
        placeholder={placeholder}
        value={current}
        onFocus={() => {
          setOpen(true);
          if (items.length && highlightedIndex === null) {
            setHighlightedIndex(0);
          }
        }}
        onChange={(e) => {
          const v = e.target.value;
          onValueChange(v); // 👈 delegate full control to parent
          setOpen(true);
          setHighlightedIndex(null);
        }}
        onKeyDown={(e) => {
          if (e.key === "ArrowDown") {
            e.preventDefault();
            moveHighlight(1);
          } else if (e.key === "ArrowUp") {
            e.preventDefault();
            moveHighlight(-1);
          } else if (e.key === "Enter") {
            if (!open) return;
            e.preventDefault();

            if (highlightedIndex !== null && items[highlightedIndex]) {
              const item = items[highlightedIndex];
              if (item.type === "option") {
                handleSelect(item.option);
              } else if (item.type === "create") {
                handleCreate();
              }
            } else if (canCreate) {
              handleCreate();
            }
          } else if (e.key === "Escape") {
            setOpen(false);
            setHighlightedIndex(null);
          }
        }}
      />

      {open && items.length > 0 && (
        <div className="absolute z-50 mt-1 max-h-60 w-full overflow-auto rounded-md border bg-popover text-popover-foreground shadow-md">
          <ul ref={listRef} className="p-1 text-sm">
            {items.map((item, index) => {
              const isHighlighted = index === highlightedIndex;

              if (item.type === "option") {
                const opt = item.option;
                const isSelected = opt.label === current;

                return (
                  <li
                    key={opt.value}
                    className={cn(
                      "flex cursor-pointer items-center justify-between rounded-sm px-2 py-1.5 hover:bg-muted",
                      isHighlighted && "bg-muted"
                    )}
                    onMouseDown={(e) => {
                      e.preventDefault(); // prevent blur before click
                      handleSelect(opt);
                    }}
                    onMouseEnter={() => setHighlightedIndex(index)}
                  >
                    <span>{opt.label}</span>
                    {isSelected && <CheckIcon className="h-4 w-4 opacity-70" />}
                  </li>
                );
              }

              return (
                <li
                  key="__create"
                  className={cn(
                    "mt-1 flex cursor-pointer items-center rounded-sm px-2 py-1.5 text-xs text-muted-foreground hover:bg-muted",
                    isHighlighted && "bg-muted"
                  )}
                  onMouseDown={(e) => {
                    e.preventDefault();
                    handleCreate();
                  }}
                  onMouseEnter={() => setHighlightedIndex(index)}
                >
                  Create “{item.label}”
                </li>
              );
            })}
          </ul>
        </div>
      )}
    </div>
  );
}
