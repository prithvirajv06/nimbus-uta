import {
  Component,
  ElementRef,
  ViewChild,
  AfterViewInit,
  ChangeDetectionStrategy,
  NgZone,
  Input
} from '@angular/core';

@Component({
  selector: 'app-pan-zoom-canvas',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div
      #viewport
      class="relative h-full w-full overflow-hidden
             bg-background cursor-grab active:cursor-grabbing
             touch-none select-none"
      (pointerdown)="onPointerDown($event)"
      (pointermove)="onPointerMove($event)"
      (pointerup)="onPointerUp()"
      (pointerleave)="onPointerUp()"
      (wheel)="onWheel($event)"
    >
      <div
        #content
        class="relative origin-top-left transform-gpu
               will-change-transform"
        style="width:fit-content; height:fit-content;"
      >
        <ng-content></ng-content>
      </div>
    </div>
  `
})
export class PanZoomCanvasComponent implements AfterViewInit {

  /* ---------- CONFIG ---------- */

  @Input() canvasWidth = 200;
  @Input() canvasHeight = 200;

  @Input() minScale = 0.25;
  @Input() maxScale = 2.5;

  @ViewChild('viewport', { static: true })
  viewportRef!: ElementRef<HTMLDivElement>;

  @ViewChild('content', { static: true })
  contentRef!: ElementRef<HTMLDivElement>;

  /* ---------- TRANSFORM STATE ---------- */

  private x = 0;
  private y = 0;
  private scale = 1;

  /* ---------- PAN STATE ---------- */

  private isPanning = false;
  private startX = 0;
  private startY = 0;

  /* ---------- PERFORMANCE ---------- */

  private raf = 0;
  private zoomDelta = 0;
  private zooming = false;

  constructor(private zone: NgZone) {}

  ngAfterViewInit() {
    this.zone.runOutsideAngular(() => {
      this.fitToScreen();
    });
  }

  /* ============================================================
     PAN — ONLY when clicking EMPTY CANVAS
     ============================================================ */

  onPointerDown(e: PointerEvent) {
    // 🔒 CRITICAL FIX:
    // Do NOT pan when clicking child elements
    if (e.target !== this.viewportRef.nativeElement) return;

    if (e.button !== 0) return;

    this.isPanning = true;
    this.startX = e.clientX - this.x;
    this.startY = e.clientY - this.y;

    this.viewportRef.nativeElement
      .setPointerCapture(e.pointerId);
  }

  onPointerMove(e: PointerEvent) {
    if (!this.isPanning) return;

    this.x = e.clientX - this.startX;
    this.y = e.clientY - this.startY;

    this.scheduleRender();
  }

  onPointerUp() {
    this.isPanning = false;
  }

  /* ============================================================
     ZOOM — Smooth, GPU, Non-jitter
     ============================================================ */

  onWheel(e: WheelEvent) {
    e.preventDefault();

    this.zoomDelta += -e.deltaY * 0.002;

    if (this.zooming) return;
    this.zooming = true;

    requestAnimationFrame(() => {
      this.applyZoom(e);
      this.zoomDelta = 0;
      this.zooming = false;
    });
  }

  private applyZoom(e: WheelEvent) {
    const prevScale = this.scale;

    let nextScale = prevScale * (1 + this.zoomDelta);
    nextScale = Math.round(nextScale * 1000) / 1000;
    nextScale = this.clamp(nextScale);

    const rect = this.viewportRef.nativeElement.getBoundingClientRect();
    const mx = e.clientX - rect.left;
    const my = e.clientY - rect.top;

    this.x -= (mx - this.x) * (nextScale / prevScale - 1);
    this.y -= (my - this.y) * (nextScale / prevScale - 1);

    this.scale = nextScale;
    this.render();
  }

  /* ============================================================
     FIT TO SCREEN
     ============================================================ */

  fitToScreen() {
    const viewport = this.viewportRef.nativeElement.getBoundingClientRect();

    const scaleX = viewport.width / this.canvasWidth;
    const scaleY = viewport.height / this.canvasHeight;

    this.scale = this.clamp(Math.min(scaleX, scaleY));

    this.x = (viewport.width - this.canvasWidth * this.scale) / 2;
    this.y = (viewport.height - this.canvasHeight * this.scale) / 2;

    this.render();
  }

  /* ============================================================
     RENDER (FAST PATH)
     ============================================================ */

  private scheduleRender() {
    cancelAnimationFrame(this.raf);
    this.raf = requestAnimationFrame(() => this.render());
  }

  private render() {
    this.contentRef.nativeElement.style.transform =
      `translate3d(${this.x}px, ${this.y}px, 0)
       scale3d(${this.scale}, ${this.scale}, 1)`;
  }

  private clamp(v: number) {
    return Math.min(this.maxScale, Math.max(this.minScale, v));
  }
}
