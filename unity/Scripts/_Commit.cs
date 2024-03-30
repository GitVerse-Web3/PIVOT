using System;
using System.Collections;
using System.Collections.Generic;
using ModestTree;
using Zenject;

public class _Commit : ICommit
{
	public class Factory : IFactory<int, ICommit>
	{
		Random random;

		public Factory()
		{
			random = new Random();
		}

		public ICommit Create(int param)
		{
			_Commit ans = new _Commit(
				(long)random.Next() * (long)random.Next(),





		}
	}


	protected _Commit(long modelHashID, PublicKey author, long authorSignature, DateTime timestamp, Tag tag, string commitMessage, double compressionRatio, ICommit parentModel)
	{
		this.modelHashID = modelHashID;
		this.author = author;
		this.authorSignature = authorSignature;
		this.timestamp = timestamp;
		this.tag = tag;
		this.commitMessage = commitMessage;
		this.compressionRatio = compressionRatio;
		this.parentModel = parentModel;
	}

	public Int64 modelHashID { get; }
	public PublicKey author { get; }
	public Int64 authorSignature { get; }
	public DateTime timestamp { get; }
	public Tag tag { get; }
	public string commitMessage { get; }

	public double compressionRatio { get; }

	public ICommit parentModel { get; }

	public virtual bool checkValid()
	{
		return true;
	}

	public virtual byte[] getFullModel()
	{
		return new byte[1];
	}

	public virtual void rebaseToMaster()
	{
		Assert.That(this.parentModel.tag.isHead || this.parentModel == null);
		this.tag.isMaster = true;
		this.parentModel.tag.isHead = false;
		this.tag.isHead = true;
	}
}
