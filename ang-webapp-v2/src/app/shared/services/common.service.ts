import { HttpClient } from "@angular/common/http";
import { inject } from "@angular/core";
import { NotificationService } from "./notification.service";

export class CommonService {
    notificationService = inject(NotificationService)
    httpClient = inject(HttpClient);
    baseApiUrl = "http://localhost:8080/api/v1";
}