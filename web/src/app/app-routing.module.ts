import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { AppInfoComponent } from './app-info/app-info.component';
import { AppListComponent } from './app-list/app-list.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { NewAppComponent } from './new-app/new-app.component';
import { NewUpdateComponent } from './new-update/new-update.component';
import { LoginComponent } from './login/login.component';
import { RegisterComponent } from './register/register.component';
import { ReviewComponent } from './review/review.component';
import { PublishComponent } from './publish/publish.component';
import { LandingComponent } from './landing/landing.component';
import { ConsoleLayoutComponent } from './console-layout/console-layout.component';
import { UnauthorizedRegisterComponent } from './unauthorized-register/unauthorized-register.component';
import { AuthGuard } from './auth.guard';
import { ReviewerGuard } from './reviewer.guard';
import { PublisherGuard } from './publisher.guard';

const routes: Routes = [
    { path: '', component: LandingComponent },
    { path: 'auth/github/callback', component: LoginComponent },
    { path: 'register', component: RegisterComponent },
    { path: '', component: ConsoleLayoutComponent, canActivate: [AuthGuard], children: [
        { path: 'dashboard', component: DashboardComponent },
        { path: 'apps', component: AppListComponent, },
        { path: 'apps/:id', component: AppInfoComponent },
        { path: 'apps/:id/update', component: NewUpdateComponent },
        { path: 'new-app', component: NewAppComponent },
        { path: 'review', component: ReviewComponent, canActivate: [ReviewerGuard] },
        { path: 'publish', component: PublishComponent, canActivate: [PublisherGuard] },
    ] },
    { path: 'unauthorized-register', component: UnauthorizedRegisterComponent },
];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
})
export class AppRoutingModule { }
