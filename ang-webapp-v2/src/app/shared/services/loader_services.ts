import { Injectable, signal } from '@angular/core';

@Injectable({ providedIn: 'root' })
export class LoadingService {
  // A signal to track the loading state
  private _loading = signal<boolean>(false);

  // Public read-only signal
  public isLoading = this._loading.asReadonly();

  show() {
    this._loading.set(true);
  }

  hide() {
    // Small delay to ensure smooth transition
    setTimeout(() => this._loading.set(false), 500);
  }
}
