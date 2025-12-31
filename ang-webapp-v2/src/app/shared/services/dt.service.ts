import { Injectable } from "@angular/core";
import { DecisionTable } from "../types/dt.type";
import { Observable, catchError } from "rxjs";
import { ApiResponse } from "../types/common.type";
import { VariablePackage } from "../types/variable_package";
import { CommonService } from "./common.service";


@Injectable({
    providedIn: 'root'
})
export class DtService extends CommonService {

    getList(data: any): Observable<ApiResponse<DecisionTable[]>> {
        return this.httpClient.post<ApiResponse<DecisionTable[]>>(`${this.baseApiUrl}/decision-tables/list`, data)
            .pipe(catchError((error: any) => { this.handleError(error); throw error; }));
    }

    get(id: string, version: number) {
        return this.httpClient.get<ApiResponse<DecisionTable>>(`${this.baseApiUrl}/decision-table`, { params: { nimb_id: id, version: version } })
            .pipe(catchError((error: any) => { this.handleError(error); throw error; }));
    }

    create(data: any) {
        return this.httpClient.post<ApiResponse<DecisionTable>>(`${this.baseApiUrl}/decision-tables`, data)
            .pipe(catchError((error: any) => { this.handleError(error); throw error; }));
    }

    update(id: string, version: number, data: any) {
        return this.httpClient.put<ApiResponse<DecisionTable>>(`${this.baseApiUrl}/decision-tables`, data, { params: { nimb_id: id, version: version } })
            .pipe(catchError((error: any) => { this.handleError(error); throw error; }));
    }

    delete(id: string, version: number) {
        return this.httpClient.delete<ApiResponse<DecisionTable>>(`${this.baseApiUrl}/decision-table`, { params: { nimb_id: id, version: version } })
            .pipe(catchError((error: any) => { this.handleError(error); throw error; }));
    }
}