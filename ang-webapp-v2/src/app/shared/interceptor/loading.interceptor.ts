import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { finalize } from 'rxjs';
import { LoadingService } from '../services/loader_services';

export const loadingInterceptor: HttpInterceptorFn = (req, next) => {
  const loadingService = inject(LoadingService);
  if (req.url.includes('ai-worker')) {
    return next(req);
  }
  // Turn on loader
  loadingService.show();

  return next(req).pipe(
    // Turn off loader when request finishes (success or error)
    finalize(() => loadingService.hide())
  );
};
