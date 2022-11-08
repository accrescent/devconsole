import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';

import { Observable } from 'rxjs';

import { LoginResult } from './login-result';

@Injectable({
    providedIn: 'root'
})
export class AuthService {
    private readonly authCallbackUrl = 'api/auth/github/callback';
    private readonly sessionUrl = 'api/session';

    constructor(private http: HttpClient) {}

    logIn(code: string, state: string): Observable<LoginResult> {
        const params = new HttpParams().append('code', code).append('state', state);

        return this.http.get<LoginResult>(this.authCallbackUrl, { params });
    }

    logOut(): Observable<void> {
        localStorage.removeItem('loggedIn');
        localStorage.removeItem('registered');
        localStorage.removeItem('reviewer');
        localStorage.removeItem('publisher');

        return this.http.delete<void>(this.sessionUrl);
    }

    get loggedIn(): boolean {
        return localStorage.getItem('loggedIn') === 'true';
    }

    set loggedIn(loggedIn: boolean) {
        localStorage.setItem('loggedIn', String(loggedIn));
    }

    get registered(): boolean {
        return localStorage.getItem('registered') === 'true';
    }

    set registered(registered: boolean) {
        localStorage.setItem('registered', String(registered));
    }

    get reviewer(): boolean {
        return localStorage.getItem('reviewer') === 'true';
    }

    set reviewer(reviewer: boolean) {
        localStorage.setItem('reviewer', String(reviewer));
    }

    get publisher(): boolean {
        return localStorage.getItem('publisher') === 'true';
    }

    set publisher(publisher: boolean) {
        localStorage.setItem('publisher', String(publisher));
    }
}
