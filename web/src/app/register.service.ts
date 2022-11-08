import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { Observable } from 'rxjs';

import { Email } from './email';

@Injectable({
    providedIn: 'root'
})
export class RegisterService {
    private readonly registerUrl = 'api/register';
    private readonly emailsUrl = 'api/emails';

    constructor(private http: HttpClient) {}

    getEmails(): Observable<string[]> {
        return this.http.get<string[]>(this.emailsUrl);
    }

    register(email: Email): Observable<Email> {
        return this.http.post<Email>(this.registerUrl, email);
    }
}
