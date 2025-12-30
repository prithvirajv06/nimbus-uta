import { NgClass, DatePipe } from "@angular/common";
import { Component, inject } from "@angular/core";
import { PageBreadcrumbComponent } from "../../../../shared/components/common/page-breadcrumb/page-breadcrumb.component";
import { BasicTableComponent } from "../../../../shared/components/tables/basic-tables/basic-table.component";
import { ButtonComponent } from "../../../../shared/components/ui/button/button.component";
import { ModalComponent } from "../../../../shared/components/ui/modal/modal.component";
import { DecisionTable } from "../../../../shared/types/dt.type";
import { DtService } from "../../../../shared/services/dt.service";

@Component({
  selector: 'app-dt-table',
  imports: [PageBreadcrumbComponent, NgClass, ButtonComponent,
    ModalComponent,
    DatePipe],
  templateUrl: './dt-table.component.html',
  styleUrl: './dt-table.component.css',
})
export class DtTableComponent extends BasicTableComponent<DecisionTable> {

  override setService(): void {
    this.service = inject(DtService)
  }
}
