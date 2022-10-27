import { ComponentFixture, TestBed } from '@angular/core/testing';

import { NewAppComponent } from './new-app.component';

describe('NewAppComponent', () => {
    let component: NewAppComponent;
    let fixture: ComponentFixture<NewAppComponent>;

    beforeEach(async () => {
        await TestBed.configureTestingModule({
            declarations: [ NewAppComponent ]
        })
            .compileComponents();

        fixture = TestBed.createComponent(NewAppComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
