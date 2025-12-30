import { Routes } from '@angular/router';
import { DashboardComponent } from './pages/dashboard/dashboard.component';
import { NotFoundComponent } from './pages/other-page/not-found/not-found.component';
import { AppLayoutComponent } from './shared/layout/app-layout/app-layout.component';
import { SignInComponent } from './pages/auth-pages/sign-in/sign-in.component';
import { SignUpComponent } from './pages/auth-pages/sign-up/sign-up.component';
import { VariablePackageLayoutComponent } from './pages/variable-package-layout/variable-package-layout.component';
import { DecisionTableLayoutComponent } from './pages/decision-table-layout/decision-table-layout.component';

export const routes: Routes = [
  {
    path: 'app',
    component: AppLayoutComponent,
    children: [
      {
        path: 'dashboard',
        component: DashboardComponent,
        title: 'Dashboard'
      },
      {
        path: 'variable-packages',
        component: VariablePackageLayoutComponent,
        title: 'Variable Packages'
      }, {
        path: 'decision-tables',
        component: DecisionTableLayoutComponent
      }
    ],
  },
  { path: '', redirectTo: 'auth/sign-in', pathMatch: 'full' },
  {
    path: 'auth/sign-in',
    component: SignInComponent,
    title: 'Sign In Dashboard'
  },
  {
    path: 'auth/sign-up',
    component: SignUpComponent,
    title: 'Sign Up Dashboard'
  },
  {
    path: '**',
    component: NotFoundComponent,
    title: 'Not Found'
  },
];



// children:[
//   {
//     path: '',
//     component: EcommerceComponent,
//     pathMatch: 'full',
//     title:
//       'Angular Ecommerce Dashboard | TailAdmin - Angular Admin Dashboard Template',
//   },
//   {
//     path:'calendar',
//     component:CalenderComponent,
//     title:'Angular Calender | TailAdmin - Angular Admin Dashboard Template'
//   },
//   {
//     path:'profile',
//     component:ProfileComponent,
//     title:'Angular Profile Dashboard | TailAdmin - Angular Admin Dashboard Template'
//   },
//   {
//     path:'form-elements',
//     component:FormElementsComponent,
//     title:'Angular Form Elements Dashboard | TailAdmin - Angular Admin Dashboard Template'
//   },
//   {
//     path:'basic-tables',
//     component:BasicTablesComponent,
//     title:'Angular Basic Tables Dashboard | TailAdmin - Angular Admin Dashboard Template'
//   },
//   {
//     path:'blank',
//     component:BlankComponent,
//     title:'Angular Blank Dashboard | TailAdmin - Angular Admin Dashboard Template'
//   },
//   // support tickets
//   {
//     path:'invoice',
//     component:InvoicesComponent,
//     title:'Angular Invoice Details Dashboard | TailAdmin - Angular Admin Dashboard Template'
//   },
//   {
//     path:'line-chart',
//     component:LineChartComponent,
//     title:'Angular Line Chart Dashboard | TailAdmin - Angular Admin Dashboard Template'
//   },
//   {
//     path:'bar-chart',
//     component:BarChartComponent,
//     title:'Angular Bar Chart Dashboard | TailAdmin - Angular Admin Dashboard Template'
//   },
//   {
//     path:'alerts',
//     component:AlertsComponent,
//     title:'Angular Alerts Dashboard | TailAdmin - Angular Admin Dashboard Template'
//   },
//   {
//     path:'avatars',
//     component:AvatarElementComponent,
//     title:'Angular Avatars Dashboard | TailAdmin - Angular Admin Dashboard Template'
//   },
//   {
//     path:'badge',
//     component:BadgesComponent,
//     title:'Angular Badges Dashboard | TailAdmin - Angular Admin Dashboard Template'
//   },
//   {
//     path:'buttons',
//     component:ButtonsComponent,
//     title:'Angular Buttons Dashboard | TailAdmin - Angular Admin Dashboard Template'
//   },
//   {
//     path:'images',
//     component:ImagesComponent,
//     title:'Angular Images Dashboard | TailAdmin - Angular Admin Dashboard Template'
//   },
//   {
//     path:'videos',
//     component:VideosComponent,
//     title:'Angular Videos Dashboard | TailAdmin - Angular Admin Dashboard Template'
//   },
// ]