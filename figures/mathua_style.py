"""Mathua publication-figure style adapter.

Clean-room implementation of the figures4papers `scientific-figure-making`
skill contracts (`SKILL.md`, `references/api.md`), carrying Mathua design
tokens (`DESIGN.md`) instead of the upstream house style. Nothing here is
copied from upstream; only the function shapes and layout/export
conventions are mirrored so the skill's patterns apply directly.

Semantic color mapping (DESIGN.md light-mode tokens; figures print on
parchment, so dark-mode tokens never apply):
  gold    = key result / proposed method (takes the skill's blue_main slot)
  green   = improvements / positive variants
  red     = baselines / contrasts / alternatives
  teal    = decay / review-related series (extra-series slot)
  blue    = secondary series (links/interactive accent)
  neutral = reference / background categories

Typography: Space Grotesk for titles/labels, JetBrains Mono for ticks,
bar annotations and legends (matches DESIGN.md "Metric Numbers" role).
The fetchable Space Grotesk static file registers under the family name
"Space Grotesk Light" (Fontsource packaging quirk); it is listed first
in the stack so local PNG/PDF exports use real glyphs. SVG exports keep
text as <text> elements, so the family names resolve from the site's own
webfonts in the browser regardless.

Omitted from the skill's surface on purpose:
  make_sphere_illustration -- the skill itself excludes dominant-3D work;
  Mathua's 3D graph (MathConceptGraph3D) is a separate WebGL pipeline.
"""

from __future__ import annotations

import logging
import warnings
from dataclasses import dataclass, field
from pathlib import Path

import matplotlib

matplotlib.use("Agg")  # headless/batch-safe; skill convention (api.md)

import matplotlib.pyplot as plt
import numpy as np
from matplotlib import font_manager
from matplotlib.axes import Axes
from matplotlib.container import BarContainer
from matplotlib.figure import Figure

# --------------------------------------------------------------------------
# Tokens (DESIGN.md light mode; figure face is parchment, never dark #0b0f1a)
# --------------------------------------------------------------------------

PARCHMENT = "#fefcf4"
SURFACE = "#faf8f0"
INK = "#1e293b"
SECONDARY = "#475569"
MUTED = "#94a3b8"
BORDER = "#d6d3c8"

MATHUA_PALETTE = {
    "gold": "#b8933e",  # theorem gold -- key result (skill: blue_main slot)
    "gold_soft": "#e6d5a8",  # light gold band (skill: green_1-style slot)
    "green": "#059669",  # proof green -- improvements
    "green_soft": "#bfe6d4",  # light green band
    "red": "#dc2626",  # errors -- baselines / contrasts
    "red_soft": "#f3c8c8",  # light red band
    "teal": "#3a9a8a",  # decay teal -- extra series
    "blue": "#2563eb",  # integral blue -- secondary series
    "ink": INK,  # edges, text
    "muted": MUTED,  # neutral reference categories
    "secondary": SECONDARY,
}

# Skill's DEFAULT_COLORS role order, remapped: key, improvement, contrast,
# extra, secondary, neutral.
DEFAULT_COLORS = [
    MATHUA_PALETTE["gold"],
    MATHUA_PALETTE["green"],
    MATHUA_PALETTE["red"],
    MATHUA_PALETTE["teal"],
    MATHUA_PALETTE["blue"],
    MATHUA_PALETTE["muted"],
]

TEXT_STACK = ["Space Grotesk Light", "Space Grotesk", "DejaVu Sans", "sans-serif"]
MONO_STACK = ["JetBrains Mono", "DejaVu Sans Mono", "monospace"]

SUPPORTED_FORMATS = ("pdf", "svg", "eps", "png", "jpg", "jpeg", "tif", "tiff")


@dataclass(frozen=True)
class MathuaStyle:
    font_size: int = 16  # 24 for large comparison bars, 15-16 compact
    axes_linewidth: float = 2.0  # 3 for big bars, 2 for compact figures
    use_tex: bool = False  # Mathua has no TeX step; mathtext only
    font_family: tuple[str, ...] = field(default_factory=lambda: tuple(TEXT_STACK))


def register_fonts(font_dir: str | Path | None = None) -> list[str]:
    """Register figures/fonts/*.ttf with matplotlib; warn and fall back if absent."""
    directory = Path(font_dir) if font_dir else Path(__file__).resolve().parent / "fonts"
    registered: list[str] = []
    if not directory.is_dir():
        warnings.warn(f"font dir {directory} missing (run figures/fetch_fonts.py); using fallback stack")
        return registered
    for ttf in sorted(directory.glob("*.ttf")):
        try:
            font_manager.fontManager.addfont(str(ttf))
            registered.append(ttf.name)
        except Exception as exc:  # keep rendering with fallbacks, never crash
            warnings.warn(f"could not register {ttf}: {exc}")
    if not registered:
        warnings.warn(f"no fonts in {directory} (run figures/fetch_fonts.py); using fallback stack")
    return registered


class _OnceFilter(logging.Filter):
    """Let each distinct font-manager message through once; silence repeats.

    The TEXT_STACK deliberately names "Space Grotesk" (browser-resolvable in
    SVG via the site's webfonts) ahead of locally registered files, so the
    resolver logs one "not found" per process run. That first line is the
    signal; the following thirty are noise.
    """

    def __init__(self) -> None:
        super().__init__()
        self._seen: set[str] = set()

    def filter(self, record: logging.LogRecord) -> bool:
        msg = record.getMessage()
        if msg in self._seen:
            return False
        self._seen.add(msg)
        return True


def _dedupe_font_warnings() -> None:
    logger = logging.getLogger("matplotlib.font_manager")
    if not any(isinstance(f, _OnceFilter) for f in logger.filters):
        logger.addFilter(_OnceFilter())


def apply_publication_style(style: MathuaStyle | None = None) -> MathuaStyle:
    """Configure matplotlib rcParams: parchment face, minimalist spines, vector text."""
    style = style or MathuaStyle()
    _dedupe_font_warnings()
    plt.rcParams.update(
        {
            "font.family": list(style.font_family),
            "font.size": style.font_size,
            "axes.facecolor": PARCHMENT,
            "figure.facecolor": PARCHMENT,
            "savefig.facecolor": PARCHMENT,
            "text.color": INK,
            "axes.edgecolor": INK,
            "axes.labelcolor": INK,
            "xtick.color": SECONDARY,
            "ytick.color": SECONDARY,
            "axes.spines.right": False,
            "axes.spines.top": False,
            "axes.linewidth": style.axes_linewidth,
            "axes.grid": False,
            "legend.frameon": False,
            "legend.fontsize": style.font_size - 2,
            "svg.fonttype": "none",  # editable <text> in vector exports
            "text.usetex": style.use_tex,
        }
    )
    return style


def apply_mono_ticks(ax: Axes) -> None:
    """Set tick + axis labels to the mono stack (DESIGN.md metric/labels role)."""
    for label in list(ax.get_xticklabels()) + list(ax.get_yticklabels()):
        label.set_fontfamily(MONO_STACK)
    ax.xaxis.label.set_fontfamily(MONO_STACK)
    ax.yaxis.label.set_fontfamily(MONO_STACK)


def create_subplots(nrows: int = 1, ncols: int = 1, figsize: tuple[float, float] | None = None, **kwargs) -> tuple[Figure, np.ndarray]:
    """Return (fig, axes) with axes a flattened 1D array (skill api.md contract)."""
    if nrows < 1 or ncols < 1:
        raise ValueError("nrows/ncols must be >= 1")
    fig, axes = plt.subplots(nrows, ncols, figsize=figsize, **kwargs)
    return fig, np.atleast_1d(axes).ravel()


def finalize_figure(fig: Figure, out_path: str | Path, formats: list[str] | None = None, dpi: int = 300, close: bool = True, pad: float = 2.0, **kwargs) -> list[Path]:
    """tight_layout + save to one or more formats; creates parents; returns paths."""
    base = Path(out_path)
    if formats is None:
        formats = [base.suffix.lstrip(".") or "svg"]
    for fmt in formats:
        if fmt.lower() not in SUPPORTED_FORMATS:
            raise ValueError(f"unsupported format {fmt!r}; expected one of {SUPPORTED_FORMATS}")
    fig.tight_layout(pad=pad)
    saved: list[Path] = []
    for fmt in dict.fromkeys(f.lower() for f in formats):  # dedupe, keep order
        target = base.with_suffix(f".{fmt}") if base.suffix.lstrip(".").lower() != fmt else base
        target.parent.mkdir(parents=True, exist_ok=True)
        fig.savefig(target, dpi=dpi, bbox_inches="tight")
        saved.append(target)
    if close:
        plt.close(fig)
    return saved


# --------------------------------------------------------------------------
# Validation (skill api.md rules, same error surface)
# --------------------------------------------------------------------------


def _as_arrays(series: list) -> list[np.ndarray]:
    try:
        return [np.asarray(s, dtype=float).ravel() for s in series]
    except (ValueError, TypeError) as exc:
        raise ValueError(f"series must be numeric 1D sequences: {exc}") from exc


# --------------------------------------------------------------------------
# Plot helpers
# --------------------------------------------------------------------------


def make_grouped_bar(
    ax: Axes,
    categories: list[str],
    series: list,
    labels: list[str],
    ylabel: str = "Value",
    colors: list[str] | None = None,
    annotate: bool = False,
    hatches: list[str] | None = None,
    width: float = 0.8,
) -> BarContainer:
    """Grouped bars with print-safe ink edges; optional hatch for grayscale."""
    if len(labels) != len(series):
        raise ValueError(f"{len(labels)} labels for {len(series)} series")
    arrs = _as_arrays(series)
    if any(a.shape[0] != len(categories) for a in arrs):
        raise ValueError("each series must match len(categories)")
    colors = colors or DEFAULT_COLORS
    n = len(arrs)
    x = np.arange(len(categories))
    bar_w = width / n
    last: BarContainer | None = None
    for i, (arr, lab) in enumerate(zip(arrs, labels)):
        hatch = hatches[i] if hatches and i < len(hatches) else None
        last = ax.bar(
            x + (i - (n - 1) / 2) * bar_w,
            arr,
            width=bar_w * 0.92,
            label=lab,
            color=colors[i % len(colors)],
            edgecolor=INK,  # strong edge treatment (skill common-patterns)
            linewidth=1.5,
            hatch=hatch,
            zorder=3,
        )
    ax.set_xticks(x)
    ax.set_xticklabels(categories)
    ax.set_ylabel(ylabel)
    apply_mono_ticks(ax)
    if annotate and last is not None:
        annotate_bars(ax, last)
    return last  # type: ignore[return-value]


def annotate_bars(ax: Axes, bars: BarContainer, fmt: str = "{:.2f}", fontsize: int = 10, padding: int = 3) -> None:
    """Print values above each bar (skill: in-place annotation, mono per Mathua)."""
    for bar in bars:
        height = bar.get_height()
        ax.text(
            bar.get_x() + bar.get_width() / 2,
            height + padding / 72 * (ax.get_ylim()[1] - ax.get_ylim()[0]) / 10,
            fmt.format(height),
            ha="center",
            va="bottom",
            fontsize=fontsize,
            fontfamily=MONO_STACK,
            color=INK,
        )


def make_trend(
    ax: Axes,
    x,
    y_series: list,
    labels: list[str],
    colors: list[str] | None = None,
    ylabel: str | None = None,
    xlabel: str | None = None,
    show_shadow: bool = True,
    spreads: list | None = None,
    linewidth: float = 2.5,
) -> None:
    """2-4 lines with optional uncertainty bands (skill trend encoding)."""
    x = np.asarray(x, dtype=float).ravel()
    arrs = _as_arrays(y_series)
    if len(labels) != len(arrs):
        raise ValueError(f"{len(labels)} labels for {len(arrs)} series")
    if any(a.shape[0] != x.shape[0] for a in arrs):
        raise ValueError("each series must match len(x)")
    if len(arrs) > 4:
        warnings.warn("more than 4 curves harms readability (skill trend rule)")
    colors = colors or DEFAULT_COLORS
    for i, (arr, lab) in enumerate(zip(arrs, labels)):
        color = colors[i % len(colors)]
        (line,) = ax.plot(x, arr, label=lab, color=color, linewidth=linewidth, zorder=3)
        if show_shadow:
            if spreads is not None:
                delta = np.asarray(spreads[i], dtype=float).ravel()
            else:
                delta = np.std(arr) / 2 * np.ones_like(arr)
            ax.fill_between(x, arr - delta, arr + delta, color=color, alpha=0.15, linewidth=0)
            line.set_alpha(0.95)
    if ylabel:
        ax.set_ylabel(ylabel)
    if xlabel:
        ax.set_xlabel(xlabel)
    apply_mono_ticks(ax)


def make_heatmap(
    ax: Axes,
    matrix,
    x_labels: list[str] | None = None,
    y_labels: list[str] | None = None,
    cmap: str = "YlOrBr",
    cbar_label: str | None = None,
    annotate: bool = False,
) -> None:
    """2D heatmap; default cmap is a gold ramp matching the Mathua palette."""
    mat = np.asarray(matrix, dtype=float)
    if mat.ndim != 2:
        raise ValueError("matrix must be 2D")
    im = ax.imshow(mat, cmap=cmap, aspect="auto")
    if x_labels is not None:
        ax.set_xticks(range(len(x_labels)), x_labels)
    if y_labels is not None:
        ax.set_yticks(range(len(y_labels)), y_labels)
    for label in list(ax.get_xticklabels()) + list(ax.get_yticklabels()):
        label.set_fontfamily(MONO_STACK)
    if cbar_label is not None:
        cbar = ax.figure.colorbar(im, ax=ax)
        cbar.set_label(cbar_label)
    if annotate:
        for r in range(mat.shape[0]):
            for c in range(mat.shape[1]):
                ax.text(c, r, f"{mat[r, c]:.2f}", ha="center", va="center", fontsize=9, fontfamily=MONO_STACK, color=INK)


def make_scatter(ax: Axes, x, y, label: str | None = None, color: str | None = None, size: float = 50, alpha: float = 0.7) -> None:
    """Single-series scatter (skill api.md contract)."""
    xs = np.asarray(x, dtype=float).ravel()
    ys = np.asarray(y, dtype=float).ravel()
    if xs.shape != ys.shape:
        raise ValueError("x and y must have the same length")
    ax.scatter(xs, ys, label=label, color=color or MATHUA_PALETTE["gold"], s=size, alpha=alpha, edgecolors=INK, linewidths=0.5, zorder=3)
    apply_mono_ticks(ax)
