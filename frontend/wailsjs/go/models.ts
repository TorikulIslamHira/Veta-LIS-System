export namespace core {
	
	export class LabResult {
	    id: number;
	    pet_id: number;
	    machine_name: string;
	    test_name: string;
	    result_value: string;
	    unit: string;
	    status: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new LabResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.pet_id = source["pet_id"];
	        this.machine_name = source["machine_name"];
	        this.test_name = source["test_name"];
	        this.result_value = source["result_value"];
	        this.unit = source["unit"];
	        this.status = source["status"];
	        this.created_at = this.convertValues(source["created_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Pet {
	    id: number;
	    owner_id: number;
	    name: string;
	    species: string;
	    breed: string;
	    age: number;
	    gender: string;
	
	    static createFrom(source: any = {}) {
	        return new Pet(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.owner_id = source["owner_id"];
	        this.name = source["name"];
	        this.species = source["species"];
	        this.breed = source["breed"];
	        this.age = source["age"];
	        this.gender = source["gender"];
	    }
	}

}

