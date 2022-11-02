import { Injectable } from '@angular/core';

@Injectable({
    providedIn: 'root'
})
export class AuthService {
    get loggedIn(): boolean {
        return document.cookie.split(';').some(item => item.includes('logged_in=yes'));
    }
}
