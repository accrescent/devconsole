import { ComponentFixture, TestBed } from '@angular/core/testing';

import { UnauthorizedRegisterComponent } from './unauthorized-register.component';

describe('UnauthorizedRegisterComponent', () => {
    let component: UnauthorizedRegisterComponent;
    let fixture: ComponentFixture<UnauthorizedRegisterComponent>;

    beforeEach(async () => {
        await TestBed.configureTestingModule({
            declarations: [ UnauthorizedRegisterComponent ]
        })
            .compileComponents();

        fixture = TestBed.createComponent(UnauthorizedRegisterComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
