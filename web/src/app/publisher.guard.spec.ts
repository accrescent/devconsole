import { TestBed } from '@angular/core/testing';

import { PublisherGuard } from './publisher.guard';

describe('PublisherGuard', () => {
    let guard: PublisherGuard;

    beforeEach(() => {
        TestBed.configureTestingModule({});
        guard = TestBed.inject(PublisherGuard);
    });

    it('should be created', () => {
        expect(guard).toBeTruthy();
    });
});
