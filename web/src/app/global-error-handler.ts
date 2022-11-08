import { ErrorHandler, Injectable, NgZone } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { MatSnackBar } from '@angular/material/snack-bar';

@Injectable({
    providedIn: 'root',
})
export class GlobalErrorHandler implements ErrorHandler {
    constructor(private snackbar: MatSnackBar, private zone: NgZone) {}

    handleError(error: any): void {
        if (
            error instanceof HttpErrorResponse &&
            error.error !== null &&
            Object.hasOwn(error.error, 'error')
        ) {
            this.zone.run(() => {
                this.snackbar.open(error.error.error, '', { panelClass: 'snackbar-error' });
            });
        } else {
            console.error(error);
        }
    }
}
