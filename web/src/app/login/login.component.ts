import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';

import { AuthService } from '../auth.service';

@Component({
    selector: 'app-login',
    templateUrl: './login.component.html',
    styleUrls: ['./login.component.css'],
})
export class LoginComponent implements OnInit {
    constructor(
        private authService: AuthService,
        private activatedRoute: ActivatedRoute,
        private router: Router,
    ) {}

    ngOnInit(): void {
        this.activatedRoute.queryParams.subscribe(params => {
            this.authService.logIn(params['code'], params['state']).subscribe(res => {
                this.authService.loggedIn = res.logged_in;
                this.authService.registered = res.registered;
                this.authService.reviewer = res.reviewer;
                this.authService.publisher = res.publisher;

                if (this.authService.loggedIn) {
                    if (this.authService.registered) {
                        this.router.navigate(['dashboard']);
                    } else {
                        this.router.navigate(['register']);
                    }
                } else {
                    this.router.navigate(['unauthorized-register']);
                }
            });
        });
    }
}
