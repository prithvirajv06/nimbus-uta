import { HttpClient } from "@angular/common/http";
import { inject, Injectable } from "@angular/core";
import { CommonService } from "./common.service";
import { catchError } from "rxjs";


@Injectable({
    providedIn: 'root'
})
export class SimulationService extends CommonService {

    runSimulation(data: { nimb_id: string, version: number, minor_version: number, input_data: any }, type: string) {
        return this.httpClient.post<any>(`${this.simApiUrl}/execute/${type}`, data.input_data,
            {
                params: {
                    nimb_id: data.nimb_id, version: data.version, minor_version: data.minor_version,
                    is_debug: 'YES'
                }
            }
        )
            .pipe(
                catchError((error: any) => { this.handleError(error); throw error; })
            );
    }
}