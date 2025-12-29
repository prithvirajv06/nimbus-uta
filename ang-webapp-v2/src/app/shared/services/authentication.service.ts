import { Injectable } from "@angular/core";
import { CommonService } from "./common.service";
import { AppCustomer } from "../types/user.type";
import { catchError, Observable } from "rxjs";
import { ApiResponse } from "../types/common.type";

@Injectable({
    providedIn: 'root'
})
export class AuthenticationService extends CommonService {

    signUpUser(userdata: AppCustomer): Observable<ApiResponse<AppCustomer>> {
        const url = `${this.baseApiUrl}/customers/register`;
        return this.httpClient.post<ApiResponse<AppCustomer>>(url, userdata).pipe(
            catchError((error: ApiResponse<AppCustomer>) => {
                this.notificationService.error(error.message, 5);
                throw error;
            })
        );
    }

    singInUser(email: string, password: string): Observable<ApiResponse<AppCustomer>> {
        const url = `${this.baseApiUrl}/customers/login`;
        return this.httpClient.post<ApiResponse<AppCustomer>>(url, { email, password }).pipe(
            catchError((error: ApiResponse<AppCustomer>) => {
                this.notificationService.error(error.message, 5);
                throw error;
            })
        );
    }
}