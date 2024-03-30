using System.Collections;
using System.Collections.Generic;
using System.Linq;
using UnityEngine;
using Zenject;

public class _BranchManager : BranchManager
{
	IFactory<string, ICommit, ICommit> _nodeFactory;

	ICommit _masterHead;

	public _BranchManager(IFactory<string, ICommit, ICommit> nodeFactory)
	{
		_nodeFactory = nodeFactory;
		this.masterHead = null;
		this.allNodes = new List<ICommit>();
		this.newHeadCandidates = new List<ICommit>();
	}

	public ICommit masterHead
	{
		get
		{
			return _masterHead;
		}
		set
		{

			_masterHead = value;
			_masterHead.chosenToBeHead();

		}
	}

	public IList<ICommit> allNodes { get; }

	public IList<ICommit> newHeadCandidates { get; }

	public ICommit createNewNode()
	{
		ICommit ans = _nodeFactory.Create("Fuck OpenAI", masterHead);
		allNodes.Add(ans);
		return ans;
	}

	public void Init()
	{
		ICommit root = createNewNode();
		masterHead = root;


	}

	public void newMasterHead()
	{
		ICommit ans = newHeadCandidates.First();
		foreach (var item in newHeadCandidates)
		{
			if (item.compressionRatio < ans.compressionRatio)
			{
				ans = item;
			}

		}

		masterHead = ans;
	}
}
