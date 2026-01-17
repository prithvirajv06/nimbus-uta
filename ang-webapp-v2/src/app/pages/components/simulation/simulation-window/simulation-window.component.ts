import { Component, inject, Input, OnChanges, SimpleChanges, ViewChild, input } from '@angular/core';
import { Option, SelectComponent } from "../../../../shared/components/form/select/select.component";
import { LabelComponent } from "../../../../shared/components/form/label/label.component";
import { DtService } from '../../../../shared/services/dt.service';
import { LogicFlowService } from '../../../../shared/services/logic-flow.service';
import { DecisionTable } from '../../../../shared/types/dt.type';
import { LogicFlow } from '../../../../shared/types/logic-flow.type';
import { JsonPipe, UpperCasePipe } from '@angular/common';
import { ButtonComponent } from "../../../../shared/components/ui/button/button.component";
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { VariablePackageService } from '../../../../shared/services/variable-package.service';
import { SimulationService } from '../../../../shared/services/simulation.service';
import { JsonEditorComponent, JsonEditorOptions } from 'ang-jsoneditor';

@Component({
  selector: 'app-simulation-window',
  imports: [SelectComponent, LabelComponent, UpperCasePipe, ButtonComponent, FormsModule, FormsModule, ReactiveFormsModule, JsonEditorComponent],
  templateUrl: './simulation-window.component.html',
  styleUrl: './simulation-window.component.css',
})
export class SimulationWindowComponent implements OnChanges {


  @Input() type: string = '';
  @Input() selectedSim: string = '';
  option: Option[] = [];
  dtService = inject(DtService);
  logicFlowService = inject(LogicFlowService);
  variablePackageService = inject(VariablePackageService);
  simulationService = inject(SimulationService);
  nimId: string = '';
  version: string = '';
  minorVersion: string = '';
  name: string = '';

  requestJson: any = {};
  responseJson: any = {};
  executionLog: any[] = [];

  public editorOptions: JsonEditorOptions;
  public respEditorOptions: JsonEditorOptions = new JsonEditorOptions();
  public data: any;
  // optional
  @ViewChild(JsonEditorComponent, { static: false }) editor!: JsonEditorComponent;
  formGroup: FormGroup;

  constructor() {
    this.editorOptions = new JsonEditorOptions()
    this.editorOptions.modes = ['code']; // set all allowed modes
    this.editorOptions.mode = 'code'; //set only one mode
    this.respEditorOptions.modes = ['view'];
    this.respEditorOptions.mode = 'view';
    this.formGroup = new FormGroup({
      request: new FormControl({}),
      response: new FormControl({})
    });
  }
ngOnChanges(changes: SimpleChanges): void {
    if (changes['type']) {
      this.getSimulation(this.type);
    }
    if (changes['selectedSim']) {
      this.setSimulationAttribute(this.selectedSim);
    }
  }

  parseJsonSafe(input: any): any {
    try {
      return JSON.parse(input.value);
    } catch (e) {
      return {};
    }
  }
  setSimulationAttribute(selectedSim: string) {
    var parts = selectedSim.split('~');
    this.nimId = parts[0];
    this.version = parts[1];
    this.minorVersion = parts[2];
    this.name = parts[3];
    const varPackageId = parts[4];
    const varPackageVersion = parts[5];
    this.variablePackageService.get(varPackageId, parseInt(varPackageVersion, 10)).subscribe({
      next: (res) => {
        if (res.status === 'success' && res.data) {
          this.formGroup.get('request')?.setValue(res.data.sample_json ? JSON.parse(res.data.sample_json) : {});
        }
      }
    });
  }

  getService(): DtService | LogicFlowService {
    if (this.type === 'decision-table') {
      return this.dtService;
    } else if (this.type === 'logic-flow') {
      return this.logicFlowService;
    }
    return this.dtService;
  }


  getSimulation(selectedType: string) {
    (<any>this.getService().getList({})).subscribe({
      next: (resp: any) => {
        if (resp.status === 'success' && resp.data) {
          this.option = this.extractOptions(resp.data);
        }
      }
    }
    );
    // Add logic to handle the selected simulation type
  }

  extractOptions(response: DecisionTable[] | LogicFlow[]): Option[] {
    return response.map(item => ({
      label: item.name, value: item.nimb_id + `~` + item.audit.version + "~"
        + (item.audit.minor_version !== undefined ? item.audit.minor_version : '0') + `~` + item.name + "~" + item.variable_package?.nimb_id + "~" + item.variable_package?.audit.version
    }));
  }

  loadSameJSON() {
    this.requestJson = {
      "applicant_age": 35,
      "applicant_income": 75000,
      "loan_amount": 200000,
      "credit_score": 720
    };
  }

  runSimulation() {
    const payload = {
      nimb_id: this.nimId,
      version: parseInt(this.version, 10),
      minor_version: parseInt(this.minorVersion, 10),
      input_data: this.formGroup.get('request')?.value
    };
    this.simulationService.runSimulation(payload, this.type).subscribe({
      next: (res) => {
        if (res) {
          this.responseJson = res.data ? res.data : {};
          this.formGroup.get('response')?.setValue(this.responseJson);
          this.executionLog = res.data.log ? res.data.log : [];
        }
      }
    });
  }
}
