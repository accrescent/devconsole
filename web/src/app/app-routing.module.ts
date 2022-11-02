import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { DashboardComponent } from './dashboard/dashboard.component';
import { NewAppComponent } from './new-app/new-app.component';
import { RegisterComponent } from './register/register.component';
import { LandingComponent } from './landing/landing.component';
import { ConsoleLayoutComponent } from './console-layout/console-layout.component';
import { AuthGuard } from './auth.guard';

const routes: Routes = [
    { path: '', component: LandingComponent },
    { path: 'register', component: RegisterComponent },
    { path: 'console', component: ConsoleLayoutComponent, canActivate: [AuthGuard], children: [
        { path: '', redirectTo: 'dashboard', pathMatch: 'full' },
        { path: 'dashboard', component: DashboardComponent },
        { path: 'new-app', component: NewAppComponent },
    ] },
];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
})
export class AppRoutingModule { }
