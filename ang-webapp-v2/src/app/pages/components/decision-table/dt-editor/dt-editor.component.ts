import { Component, inject, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { NotificationService } from '../../../../shared/services/notification.service';
import { DtService } from '../../../../shared/services/dt.service';
import { BreadcrumbBar } from '@amcharts/amcharts5/.internal/charts/hierarchy/BreadcrumbBar';
import { PageBreadcrumbComponent } from '../../../../shared/components/common/page-breadcrumb/page-breadcrumb.component';
import { MatExpansionModule } from "@angular/material/expansion";
import { LabelComponent } from "../../../../shared/components/form/label/label.component";
import { InputFieldComponent } from "../../../../shared/components/form/input/input-field.component";
import { TextAreaComponent } from "../../../../shared/components/form/input/text-area.component";
import { DTModelutils } from '../../../../shared/types/dt.type';
import { Field } from '@angular/forms/signals';
import { ButtonComponent } from '../../../../shared/components/ui/button/button.component';

@Component({
  selector: 'app-dt-editor',
  imports: [PageBreadcrumbComponent, MatExpansionModule, LabelComponent, InputFieldComponent, TextAreaComponent, Field, ButtonComponent],
  templateUrl: './dt-editor.component.html',
  styleUrl: './dt-editor.component.css',
})
export class DtEditorComponent implements OnInit {

  activatedRoute = inject(ActivatedRoute);
  notificationService = inject(NotificationService);
  decisionTableService = inject(DtService);
  router = inject(Router);
  formModel = DTModelutils.signalModel();
  formGroup = DTModelutils.detailsFormGroup(this.formModel);

  ngOnInit(): void {
  }

  getDecisionTable() {
  }

}
