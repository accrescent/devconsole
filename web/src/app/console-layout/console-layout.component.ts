import { Component } from '@angular/core';

import { AuthService } from '../auth.service';

@Component({
    selector: 'app-console-layout',
    templateUrl: './console-layout.component.html',
    styleUrls: ['./console-layout.component.css']
})
export class ConsoleLayoutComponent {
    constructor(private authService: AuthService) {}

    get reviewer(): boolean {
        return this.authService.reviewer;
    }
}
