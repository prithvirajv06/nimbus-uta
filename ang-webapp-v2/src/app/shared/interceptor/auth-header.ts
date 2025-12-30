import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { AuthenticationService } from '../services/authentication.service';
import { AppCustomer } from '../types/user.type';
// Assuming you have an AuthService or equivalent
// that holds the user context

/**
 * Interceptor function to add custom 'zentic-user-name' and 'zentic-org-id' headers.
 */
export const AuthHeaderInterceptor: HttpInterceptorFn = (req, next) => {
  // 1. Inject the service that holds the current user state
  let userService = inject(AuthenticationService)
  const userDeatails: AppCustomer| null = userService.getCurrentUserValue();

  // 2. Retrieve the required values from the service
  // Use a nullish coalescing operator (?? '') for safety,
  // as headers must be strings.
  const userName = (userDeatails?.fname ?? '') + (userDeatails?.lname ?? '');
  const orgId = userDeatails?.organization?.nimb_id ?? '';

  // 3. Clone the request and set the new headers
  // Headers should only be added if they have valid values,
  // but for demonstration, we will set them to ensure presence.
  const modifiedReq = req.clone({
    setHeaders: {
      'user_id': userName,
      'org_id': orgId
      // Note: Header keys are case-insensitive, but convention is lowercase/kebab-case
    }
  });

  // 4. Pass the modified request to the next handler
  return next(modifiedReq);
};
