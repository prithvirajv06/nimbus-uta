import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';
import { NotificationHostComponent } from "./shared/components/ui/notification-host/notification-host.component";

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    RouterModule,
    NotificationHostComponent
],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css',
})
export class AppComponent {
  title = 'Angular Ecommerce Dashboard | TailAdmin';
}
