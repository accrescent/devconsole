import { Injectable } from '@angular/core';
import {
    HttpErrorResponse, HttpEvent, HttpInterceptor, HttpHandler, HttpRequest,
} from '@angular/common/http';
import { Router } from '@angular/router';
import { MatSnackBar } from '@angular/material/snack-bar';

import { Observable, tap } from 'rxjs';

import { AuthService } from './auth.service';

@Injectable()
export class AuthInterceptor implements HttpInterceptor {
    constructor(
        private authService: AuthService,
        private router: Router,
        private snackbar: MatSnackBar,
    ) {}

    intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
        return next.handle(req).pipe(tap(
            () => {},
            (error: any) => {
                if (error instanceof HttpErrorResponse && error.status === 401) {
                    this.authService.logOut();
                    this.router.navigate(['/']);
                    this.snackbar.open('You must be logged in to access that resource');
                }
            }
        ));
    }
}
