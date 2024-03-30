using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public interface ICommit
{
	Int64 modelHashID { get; }
	PublicKey author { get; }
	Int64 authorSignature { get; }
	DateTime timestamp { get; }
	Tag tag { get; }
	string commitMessage { get; }
	ICommit parentModel { get; set; }




	byte[] getFullModel();
	bool checkValid();
	void chosenToBeHead();
	void rebaseToMaster(ICommit head);
	double compressionRatio { get; set; }


}




public interface Tag
{
	bool isMaster { get; set; }
	bool isDeltaModel { get; }
	bool isHead { get; set; }
	bool isEncrypted { get; }
	bool isMerge { get; }
}


public class _Tag : Tag
{
	public _Tag()
	{
		this.isMaster = false;
		this.isDeltaModel = true;
		this.isHead = false;
		this.isEncrypted = false;
		this.isMerge = false;
	}

	public bool isMaster { get; set; }
	public bool isDeltaModel { get; }
	public bool isHead { get; set; }
	public bool isEncrypted { get; }
	public bool isMerge { get; }
}

public interface MergeModel : ICommit
{
	ICommit[] mergingModels { get; }

}



public interface DeltaModel : ICommit
{
	byte[] modelDelta { get; }

}

public interface FullModel : ICommit
{


	byte[] modelBin { get; }

}




public interface PublicKey
{
	string key { get; set; }
}
