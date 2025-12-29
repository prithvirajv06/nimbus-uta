import { Injectable } from "@angular/core";
import { CommonService } from "./common.service";
import { AppCustomer } from "../types/user.type";
import { catchError, Observable } from "rxjs";
import { ApiResponse } from "../types/common.type";

@Injectable({
    providedIn: 'root'
})
export class AuthenticationService extends CommonService {

    currentUser: AppCustomer | null = null;


    signUpUser(userdata: AppCustomer): Observable<ApiResponse<AppCustomer>> {
        const url = `${this.baseApiUrl}/customers/register`;
        return this.httpClient.post<ApiResponse<AppCustomer>>(url, userdata).pipe(
            catchError((httpError: any) => {
                this.notificationService.error(httpError.error.message, 5);
                throw httpError;
            })
        );
    }

    singInUser(email: string, password: string): Observable<ApiResponse<AppCustomer>> {
        const url = `${this.baseApiUrl}/customers/login`;
        return this.httpClient.post<ApiResponse<AppCustomer>>(url, { email, password }).pipe(
            catchError((httpError: any) => {
                this.notificationService.error(httpError.error.message, 5);
                throw httpError;
            })
        );
    }


    setCurrentUser(user: AppCustomer): void {
        localStorage.setItem('currentUser', JSON.stringify(user));
    }

    getCurrentUserValue(): AppCustomer | null {
        if (!this.currentUser) {
            const userData = localStorage.getItem('currentUser');
            this.currentUser = userData ? JSON.parse(userData) : null;
        }
        return this.currentUser;
    }
}