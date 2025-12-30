import { Injectable } from "@angular/core";
import { CommonService } from "./common.service";
import { ApiResponse } from "../types/common.type";
import { VariablePackage } from "../types/variable_package";
import { catchError, Observable } from "rxjs";


@Injectable({
    providedIn: 'root'
})
export class VariablePackageService extends CommonService {

    getList(data:any): Observable<ApiResponse<VariablePackage[]>> {
        return this.httpClient.post<ApiResponse<VariablePackage[]>>(`${this.baseApiUrl}/variables-packages/list`, data)
            .pipe(catchError((error: any) => { this.handleError(error); throw error; }));
    }

    get(id: string, version: number) {
        return this.httpClient.get<ApiResponse<VariablePackage>>(`${this.baseApiUrl}/variables-packages`, { params: { nimb_id: id, version: version } })
            .pipe(catchError((error: any) => { this.handleError(error); throw error; }));
    }

    create(data: any) {
        return this.httpClient.post<ApiResponse<VariablePackage>>(`${this.baseApiUrl}/variables-packages`, data)
            .pipe(catchError((error: any) => { this.handleError(error); throw error; }));
    }

    update(id: string, version: number, data: any) {
        return this.httpClient.put<ApiResponse<VariablePackage>>(`${this.baseApiUrl}/variables-packages`, data, { params: { nimb_id: id, version: version } })
            .pipe(catchError((error: any) => { this.handleError(error); throw error; }));
    }

    delete(id: string, version: number) {
        return this.httpClient.delete<ApiResponse<VariablePackage>>(`${this.baseApiUrl}/variables-package`, { params: { nimb_id: id, version: version } })
            .pipe(catchError((error: any) => { this.handleError(error); throw error; }));
    }

}