import { Component } from '@angular/core';
import { EcommerceMetricsComponent } from '../../shared/components/ecommerce/ecommerce-metrics/ecommerce-metrics.component';
import { StatisticsChartComponent } from '../../shared/components/ecommerce/statics-chart/statics-chart.component';

@Component({
  selector: 'app-ecommerce',
  imports: [
    EcommerceMetricsComponent,
    StatisticsChartComponent
],
  templateUrl: './dashboard.component.html',
})
export class DashboardComponent {}
