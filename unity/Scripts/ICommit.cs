using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;

interface ICommit
{
	Int64 modelHash { get; }
	string commitMessage { get; }
	PublicKey author { get; }
	Int64 authorSignature { get; }
	bool isMaster { get; set; }
	byte[] getFullModel();
	bool checkValid();
	void rebaseToMaster();
	double getCompressionRatio();
}


interface DeltaModel : ICommit
{
	ICommit parentModel { get; }
	bool isEncrypted { get; }
	byte[] modelDelta { get; }

}


interface MergeModel : ICommit
{
	ICommit[] parentModels { get; }

}

public interface PublicKey
{
	string key { get; set; }
}