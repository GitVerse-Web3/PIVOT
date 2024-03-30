using System;
using System.Collections;
using System.Collections.Generic;


public interface ICommit
{
	Int64 modelHashID { get; }
	PublicKey author { get; }
	Int64 authorSignature { get; }
	DateTime timestamp { get; }
	Tag tag { get; }
	string commitMessage { get; }
	ICommit parentModel { get; }




	byte[] getFullModel();
	bool checkValid();
	void rebaseToMaster();
	double compressionRatio { get; }
}




public interface Tag
{
	bool isMaster { get; set; }
	bool isDeltaModel { get; }
	bool isHead { get; set; }
	bool isEncrypted { get; }
	bool isMerge { get; }
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