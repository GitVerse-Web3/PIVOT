using System;
using System.Collections;
using System.Collections.Generic;
using ModestTree;
using Zenject;

public class _Commit : ICommit
{
	public class Factory : IFactory<string, ICommit, ICommit>
	{
		Random random;

		public Factory()
		{
			random = new Random();
		}

		public ICommit Create(string commitMessage, ICommit parentModel)
		{


			_Commit ans = new _Commit(
				(long)random.Next() * (long)random.Next()
				, new _PublicKey()
				, random.Next()
				, DateTime.Now
				, new _Tag()
				, commitMessage
				, 1
				, parentModel
				);
			ans.makeUpC();


			return ans;

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

	public double compressionRatio { get; set; }

	public ICommit parentModel { get; set; }

	protected virtual void makeUpC()
	{
		double c;
		if (parentModel == null)
		{
			c = 1;
		}
		else
		{
			float nn = 2 / 3;
			float xx = 1 - nn;
			float p = (float)parentModel.compressionRatio;

			c = nn * p + UnityEngine.Random.Range(0.5f, 1f) * xx * p;
		}
		compressionRatio = c;
	}

	public virtual bool checkValid()
	{
		return true;
	}

	public virtual byte[] getFullModel()
	{
		return new byte[1];
	}

	public virtual void chosenToBeHead()
	{


		this.tag.isMaster = true;
		this.tag.isHead = true;
		if (this.parentModel != null)
		{
			Assert.That(this.parentModel.tag.isHead);
			this.parentModel.tag.isHead = false;
		}


	}

	public virtual void rebaseToMaster(ICommit head)
	{
		this.parentModel = head;
		this.makeUpC();
	}
}
