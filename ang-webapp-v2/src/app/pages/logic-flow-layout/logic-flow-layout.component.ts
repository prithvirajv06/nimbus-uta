import { Component } from '@angular/core';
import { DecisionTable } from '../../shared/types/dt.type';
import { CommonLayoutComponent } from '../common-management-component';
import { LogicFlow } from '../../shared/types/logic-flow.type';
import { LogicFlowTableComponent } from "../components/logic-flow/logic-flow-table/logic-flow-table.component";
import { LogicFlowEditorComponent } from '../components/logic-flow/logic-flow-editor/logic-flow-editor.component';

@Component({
  selector: 'app-logic-flow-layout',
  imports: [LogicFlowTableComponent, LogicFlowEditorComponent],
  templateUrl: './logic-flow-layout.component.html',
  styleUrl: './logic-flow-layout.component.css',
})
export class LogicFlowLayoutComponent extends CommonLayoutComponent<LogicFlow> {


}
