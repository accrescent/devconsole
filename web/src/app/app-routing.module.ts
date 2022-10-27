import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { DashboardComponent } from './dashboard/dashboard.component';
import { NewAppComponent } from './new-app/new-app.component';
import { RegisterComponent } from './register/register.component';
import { LandingComponent } from './landing/landing.component';

const routes: Routes = [
    { path: '', component: LandingComponent },
    { path: 'dashboard', component: DashboardComponent },
    { path: 'register', component: RegisterComponent },
    { path: 'new-app', component: NewAppComponent },
];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
})
export class AppRoutingModule { }
