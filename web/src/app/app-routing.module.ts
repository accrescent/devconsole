import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { DashboardComponent } from './dashboard/dashboard.component';
import { NewAppComponent } from './new-app/new-app.component';
import { LoginComponent } from './login/login.component';
import { RegisterComponent } from './register/register.component';
import { ReviewComponent } from './review/review.component';
import { LandingComponent } from './landing/landing.component';
import { ConsoleLayoutComponent } from './console-layout/console-layout.component';
import { AuthGuard } from './auth.guard';
import { ReviewerGuard } from './reviewer.guard';

const routes: Routes = [
    { path: '', component: LandingComponent },
    { path: 'auth/github/callback', component: LoginComponent },
    { path: 'register', component: RegisterComponent },
    { path: '', component: ConsoleLayoutComponent, canActivate: [AuthGuard], children: [
        { path: 'dashboard', component: DashboardComponent },
        { path: 'new-app', component: NewAppComponent },
        { path: 'review', component: ReviewComponent, canActivate: [ReviewerGuard] },
    ] },
];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
})
export class AppRoutingModule { }
