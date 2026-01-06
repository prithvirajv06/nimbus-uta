import { inject } from '@angular/core';
import { Component, OnInit } from '@angular/core';
import { PageBreadcrumbComponent } from '../../shared/components/common/page-breadcrumb/page-breadcrumb.component';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { DtService } from '../../shared/services/dt.service';
import { LogicFlowService } from '../../shared/services/logic-flow.service';
import { SimulationWindowComponent } from '../components/simulation/simulation-window/simulation-window.component';

@Component({
  selector: 'app-simulation-layout',
  imports: [PageBreadcrumbComponent, RouterLink, SimulationWindowComponent],
  templateUrl: './simulation-layout.component.html',
  styleUrl: './simulation-layout.component.css',
})
export class SimulationLayoutComponent implements OnInit {
  viewWindow = 'SELECTION';
  type = '';
  activatedRoute = inject(ActivatedRoute)



  constructor() { }


  ngOnInit(): void {
    this.activatedRoute.queryParams.subscribe(params => {
      const type = params['type'];
      if (type) {
        this.type = type;
        this.setViewWindow('SIMULATION');
      } else {
        this.setViewWindow('SELECTION');
      }
    });
  }

  setViewWindow(view: string) {
    this.viewWindow = view;
  }
}
