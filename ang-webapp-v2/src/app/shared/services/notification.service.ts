import { Injectable } from '@angular/core';
import { Subject, Observable } from 'rxjs';
import { filter } from 'rxjs/operators';

// 1. Define the TypeScript Model for a Notification
export type NotificationType = 'success' | 'warning' | 'error' | 'info';

export interface Notification {
  id: number;
  type: NotificationType;
  message: string;
  duration?: number; // Time in ms before auto-dismissal
  actionLabel?: string; // Optional label for a clickable action
}

@Injectable({
  providedIn: 'root'
})
export class NotificationService {
  // 2. Use a Subject to act as the notification stream
  // Subjects are both Observables (to read from) and Observers (to push to).
  private notificationSubject = new Subject<Notification>();

  // 3. Expose the stream as a public Observable for components to subscribe to
  public notifications$: Observable<Notification> = this.notificationSubject.asObservable();

  private nextId = 1;

  // 4. Core method to send a notification
  private sendNotification(type: NotificationType, message: string, duration: number = 5000, actionLabel?: string): void {
    const notification: Notification = {
      id: this.nextId++,
      type,
      message,
      duration,
      actionLabel,
    };
    this.notificationSubject.next(notification);
  }

  // 5. Utility methods for ease of use
  success(message: string, duration?: number, actionLabel?: string): void {
    this.sendNotification('success', message, duration, actionLabel);
  }

  error(message: string, duration?: number, actionLabel?: string): void {
    this.sendNotification('error', message, duration, actionLabel); // Errors don't auto-dismiss by default
  }

  warning(message: string, duration?: number, actionLabel?: string): void {
    this.sendNotification('warning', message, duration, actionLabel);
  }

  info(message: string, duration?: number, actionLabel?: string): void {
    this.sendNotification('info', message, duration, actionLabel);
  }
}
