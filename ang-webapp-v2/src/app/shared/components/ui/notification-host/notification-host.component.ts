import { Component, OnInit, OnDestroy, inject } from '@angular/core';
import { from, Subscription } from 'rxjs';
import { NgIf, NgFor, CommonModule } from '@angular/common'; // Import necessary modules for template
import { Notification, NotificationService, NotificationType } from '../../../services/notification.service';

@Component({
  selector: 'app-notification-host',
  templateUrl: './notification-host.component.html',
  styleUrls: ['./notification-host.component.css'],
  imports: [CommonModule]
})
export class NotificationHostComponent implements OnInit, OnDestroy {
  private notificationService = inject(NotificationService);
  private subscription!: Subscription;

  notifications: Notification[] = [];

  ngOnInit(): void {
    // Subscribe to the service stream
    this.subscription = this.notificationService.notifications$.subscribe(
      (notification: Notification) => {
        this.addNotification(notification);
      }
    );
  }

  private addNotification(notification: Notification): void {
    // Add the new notification to the array
    this.notifications.unshift(notification);

    // Set an auto-dismiss timer if a duration is specified and not zero
    if (notification.duration && notification.duration > 0) {
      notification.duration = Math.max(notification.duration * 1000, 2000); // Minimum duration of 2 seconds
      setTimeout(() => this.dismiss(notification.id), (notification.duration));
    }
  }

  // --- Helper Methods for Template Display ---

  // Retrieves an icon based on the notification type
  getIconClass(type: NotificationType): string {
    switch (type) {
      case 'success':
        return 'fa-check-circle';
      case 'error':
        return 'fa-exclamation-triangle';
      case 'warning':
        return 'fa-bell';
      case 'info':
        return 'fa-info-circle';
      default:
        return 'fa-question-circle';
    }
  }

  // Retrieves a title based on the notification type
  getTitle(type: NotificationType): string {
    switch (type) {
      case 'success':
        return 'Success!';
      case 'error':
        return 'Action Required';
      case 'warning':
        return 'Caution';
      case 'info':
        return 'Notification';
      default:
        return 'Message';
    }
  }

  // --- Action and Dismiss Methods ---

  // Method to remove a notification
  dismiss(id: number): void {
    this.notifications = this.notifications.filter(n => n.id !== id);
  }

  // Method for the optional action button
  performAction(notification: Notification): void {
    // In a real app, you would use an EventEmitter here or call another service
    console.log(`Action performed for notification ID: ${notification.id} (${notification.actionLabel})`);
    this.dismiss(notification.id);
  }

  ngOnDestroy(): void {
    this.subscription.unsubscribe();
  }
}


/**
 * this.notificationService.success('Data saved successfully!', 4000); // Auto-dismiss after 4s
 * this.notificationService.warning('Your session will expire in 5 minutes.', 7000);
 */
