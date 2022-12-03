import { NgModule, ErrorHandler } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import { ReactiveFormsModule } from '@angular/forms';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatListModule } from '@angular/material/list';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatRadioModule } from '@angular/material/radio';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatSnackBarModule, MAT_SNACK_BAR_DEFAULT_OPTIONS } from '@angular/material/snack-bar';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { RegisterComponent } from './register/register.component';
import { RegisterFormComponent } from './register-form/register-form.component';
import { LandingComponent } from './landing/landing.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { NewAppComponent } from './new-app/new-app.component';
import { NewAppFormComponent } from './new-app-form/new-app-form.component';
import { NewUpdateComponent } from './new-update/new-update.component';
import { NewUpdateFormComponent } from './new-update-form/new-update-form.component';
import { ConsoleLayoutComponent } from './console-layout/console-layout.component';
import { LoginComponent } from './login/login.component';
import { ReviewComponent } from './review/review.component';
import { AppListComponent } from './app-list/app-list.component';
import { PublishComponent } from './publish/publish.component';
import { GlobalErrorHandler } from './global-error-handler';
import { AuthInterceptor } from './auth-interceptor';

@NgModule({
    declarations: [
        AppComponent,
        RegisterComponent,
        RegisterFormComponent,
        LandingComponent,
        DashboardComponent,
        NewAppComponent,
        NewAppFormComponent,
        NewUpdateComponent,
        NewUpdateFormComponent,
        ConsoleLayoutComponent,
        LoginComponent,
        ReviewComponent,
        AppListComponent,
        PublishComponent,
    ],
    imports: [
        BrowserModule,
        HttpClientModule,
        BrowserAnimationsModule,
        MatButtonModule,
        MatCardModule,
        MatIconModule,
        MatInputModule,
        MatListModule,
        MatProgressSpinnerModule,
        MatRadioModule,
        MatSidenavModule,
        MatToolbarModule,
        MatSnackBarModule,
        ReactiveFormsModule,
        AppRoutingModule
    ],
    providers: [{
        provide: ErrorHandler,
        useClass: GlobalErrorHandler,
    }, {
        provide: HTTP_INTERCEPTORS,
        useClass: AuthInterceptor,
        multi: true,
    }, {
        provide: MAT_SNACK_BAR_DEFAULT_OPTIONS,
        useValue: { duration: 5000 },
    }],
    bootstrap: [AppComponent]
})
export class AppModule { }
