﻿using System.Collections;
using System.Collections.Generic;
using System.Linq;
using Zenject;

public interface BranchManager
{
	ICommit masterHead { get; }
	void newMasterHead(); //auto select the node with least compression ratio
	IList<ICommit> allNodes { get; }
	IList<ICommit> newHeadCandidates { get; }// node.parent == corrent head
	void Init();
	ICommit createNewNode();

}
