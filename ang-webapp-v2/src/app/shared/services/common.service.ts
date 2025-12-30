import { HttpClient } from "@angular/common/http";
import { inject } from "@angular/core";
import { NotificationService } from "./notification.service";

export class CommonService {
    notificationService = inject(NotificationService)
    httpClient = inject(HttpClient);
    baseApiUrl = "http://localhost:8080/api/v1";



    handleError(error: any) {
        if (error.error && error.error.message)
            this.notificationService.error(error.error.message, 5);
        else
            this.notificationService.error("An unexpected error occurred."+ error.message, 5);
    }
}