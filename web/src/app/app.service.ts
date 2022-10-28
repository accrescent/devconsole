import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { Observable } from 'rxjs';

import { App } from './app';

@Injectable({
    providedIn: 'root'
})
export class AppService {
    private readonly appsUrl = 'api/apps';

    constructor(private http: HttpClient) {}

    uploadApp(app: File): Observable<App> {
        const formData = new FormData();
        formData.append("app", app);

        return this.http.post<App>(this.appsUrl, formData);
    }

    submitApp(id: string): Observable<void> {
        return this.http.patch<void>(`${this.appsUrl}/${id}`, '');
    }
}
