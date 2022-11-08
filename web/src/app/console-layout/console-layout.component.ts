import { Component } from '@angular/core';
import { Router } from '@angular/router';

import { AuthService } from '../auth.service';

@Component({
    selector: 'app-console-layout',
    templateUrl: './console-layout.component.html',
    styleUrls: ['./console-layout.component.css']
})
export class ConsoleLayoutComponent {
    constructor(private authService: AuthService, private router: Router) {}

    get reviewer(): boolean {
        return this.authService.reviewer;
    }

    get publisher(): boolean {
        return this.authService.publisher;
    }

    logOut(): void {
        this.authService.logOut().subscribe();
        this.router.navigate(['/']);
    }
}
