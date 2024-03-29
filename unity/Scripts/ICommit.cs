using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;

interface ICommit
{
	Int64 modelHashID { get; }
	PublicKey author { get; }
	Int64 authorSignature { get; }
	DateTime timestamp { get; }
	Tag tag { get; }
	string commitMessage { get; }





	byte[] getFullModel();
	bool checkValid();
	void rebaseToMaster();
	double getCompressionRatio();
}

interface Tag
{
	bool isMaster { get; set; }
	bool isDeltaModel { get; }
	bool isHead { get; set; }
	bool isEncrypted { get; }
	bool isMerge { get; }
}

interface MergeModel : ICommit
{
	ICommit[] parentModels { get; }

}
interface SingleParentModel : ICommit
{
	ICommit parentModel { get; }
}


interface DeltaModel : SingleParentModel
{
	byte[] modelDelta { get; }

}

interface FullModel : SingleParentModel
{


	byte[] modelBin { get; }

}




public interface PublicKey
{
	string key { get; set; }
}