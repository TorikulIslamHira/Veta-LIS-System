export namespace drivers {
	
	export class LabResult {
	    ID: string;
	    MachineType: string;
	    TestID: string;
	    PetID: string;
	    OwnerID: string;
	    TestName: string;
	    TestCode: string;
	    ResultValue: string;
	    ResultUnit: string;
	    ReferenceRange: string;
	    Status: string;
	    Timestamp: string;
	    RawData: number[];
	    Notes: string;
	
	    static createFrom(source: any = {}) {
	        return new LabResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.MachineType = source["MachineType"];
	        this.TestID = source["TestID"];
	        this.PetID = source["PetID"];
	        this.OwnerID = source["OwnerID"];
	        this.TestName = source["TestName"];
	        this.TestCode = source["TestCode"];
	        this.ResultValue = source["ResultValue"];
	        this.ResultUnit = source["ResultUnit"];
	        this.ReferenceRange = source["ReferenceRange"];
	        this.Status = source["Status"];
	        this.Timestamp = source["Timestamp"];
	        this.RawData = source["RawData"];
	        this.Notes = source["Notes"];
	    }
	}

}

